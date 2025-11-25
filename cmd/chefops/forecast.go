package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ChefChristoph/chefops/internal"
)

// Simple helper struct for dish specs
type dishForecast struct {
	RecipeID  int
	Name      string
	Portions  float64
	YieldQty  float64
	YieldUnit string
}

// ------------------------------------------------------------
// forecast command
//
// Example:
//
//	chefops forecast \
//	  --out f1_forecast.csv \
//	  "DISH Pole Position Burger=600" \
//	  "DISH Lobster Roll=500" \
//	  "DISH Turbo Hammour Popcorn=700"
//
// ------------------------------------------------------------
func forecastCommand(args []string) {
	fs := flag.NewFlagSet("forecast", flag.ExitOnError)
	outFile := fs.String("out", "forecast.csv", "output CSV file")
	fs.Parse(args)

	specs := fs.Args()
	if len(specs) == 0 {
		fmt.Println("usage:")
		fmt.Println("  chefops forecast --out forecast.csv \"DISH Name=PORTIONS\" ...")
		fmt.Println("")
		fmt.Println("example:")
		fmt.Println("  chefops forecast --out f1_forecast.csv \\")
		fmt.Println("    \"DISH Pole Position Burger=600\" \\")
		fmt.Println("    \"DISH Lobster Roll=500\" \\")
		fmt.Println("    \"DISH Turbo Hammour Popcorn=700\"")
		os.Exit(1)
	}

	db, _ := internal.OpenDB()
	defer db.Close()

	// 1) Parse dish specs and resolve recipes
	var dishes []dishForecast

	for _, spec := range specs {
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "invalid spec (expected NAME=QTY): %s\n", spec)
			os.Exit(1)
		}

		rawName := strings.TrimSpace(parts[0])
		qtyStr := strings.TrimSpace(parts[1])

		if rawName == "" {
			fmt.Fprintf(os.Stderr, "empty dish name in spec: %s\n", spec)
			os.Exit(1)
		}

		portions, err := strconv.ParseFloat(qtyStr, 64)
		if err != nil || portions <= 0 {
			fmt.Fprintf(os.Stderr, "invalid portions in spec (need positive number): %s\n", spec)
			os.Exit(1)
		}

		recipeID, recipeName, err := findRecipeByName(db, rawName)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Fprintf(os.Stderr, "recipe not found for forecast: %s\n", rawName)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "error resolving recipe %q: %v\n", rawName, err)
			os.Exit(1)
		}

		var yieldQty float64
		var yieldUnit string
		err = db.QueryRow(`
			SELECT yield_qty, yield_unit
			FROM recipes
			WHERE id = ?
		`, recipeID).Scan(&yieldQty, &yieldUnit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading base yield for %s: %v\n", recipeName, err)
			os.Exit(1)
		}
		if yieldQty <= 0 {
			// Safety net – all your current recipes use 1.00 anyway
			yieldQty = 1
		}

		dishes = append(dishes, dishForecast{
			RecipeID:  recipeID,
			Name:      recipeName,
			Portions:  portions,
			YieldQty:  yieldQty,
			YieldUnit: yieldUnit,
		})
	}

	// 2) Aggregate ingredients (full marketlist) + direct subrecipes

	type ingAgg struct {
		Name        string
		Unit        string
		CostPerUnit float64
		TotalQty    float64
		TotalCost   float64
	}

	type subAgg struct {
		Name     string
		Unit     string
		TotalQty float64
	}

	ingredients := make(map[int]*ingAgg) // ingredient_id -> agg
	subrecipes := make(map[int]*subAgg)  // subrecipe_id -> agg

	for _, d := range dishes {
		scale := d.Portions / d.YieldQty

		// --- ingredients via recipe_items_expanded (recursive) ---
		ingRows, err := db.Query(`
			SELECT e.ingredient_id, i.name, i.unit, i.cost_per_unit, e.total_qty
			FROM recipe_items_expanded e
			JOIN ingredients i ON i.id = e.ingredient_id
			WHERE e.recipe_id = ?
		`, d.RecipeID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading ingredients for %s: %v\n", d.Name, err)
			os.Exit(1)
		}

		for ingRows.Next() {
			var ingID int
			var name, unit string
			var cpu, baseQty float64

			if err := ingRows.Scan(&ingID, &name, &unit, &cpu, &baseQty); err != nil {
				fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
				os.Exit(1)
			}

			qty := baseQty * scale

			agg, ok := ingredients[ingID]
			if !ok {
				agg = &ingAgg{
					Name:        name,
					Unit:        unit,
					CostPerUnit: cpu,
				}
				ingredients[ingID] = agg
			}

			agg.TotalQty += qty
			agg.TotalCost += qty * cpu
		}
		ingRows.Close()

		// --- direct subrecipes (for bulk prep planning) ---
		subRows, err := db.Query(`
			SELECT rs.subrecipe_id, r.name, rs.unit, rs.qty
			FROM recipe_subrecipes rs
			JOIN recipes r ON r.id = rs.subrecipe_id
			WHERE rs.recipe_id = ?
		`, d.RecipeID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading subrecipes for %s: %v\n", d.Name, err)
			os.Exit(1)
		}

		for subRows.Next() {
			var subID int
			var name, unit string
			var baseQty float64

			if err := subRows.Scan(&subID, &name, &unit, &baseQty); err != nil {
				fmt.Fprintf(os.Stderr, "scan error (subrecipes): %v\n", err)
				os.Exit(1)
			}

			qty := baseQty * scale

			agg, ok := subrecipes[subID]
			if !ok {
				agg = &subAgg{
					Name: name,
					Unit: unit,
				}
				subrecipes[subID] = agg
			}
			agg.TotalQty += qty
		}
		subRows.Close()
	}

	// 3) Open CSV file
	f, err := os.Create(*outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot create %s: %v\n", *outFile, err)
		os.Exit(1)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	// --- SECTION 1: Dish overview ---------------------------------
	_ = w.Write([]string{"# Dishes"})
	_ = w.Write([]string{"Dish", "Portions", "Base yield", "Scale factor"})

	for _, d := range dishes {
		scale := d.Portions / d.YieldQty
		base := fmt.Sprintf("%.3f %s", d.YieldQty, d.YieldUnit)
		_ = w.Write([]string{
			d.Name,
			fmt.Sprintf("%.3f", d.Portions),
			base,
			fmt.Sprintf("%.3f", scale),
		})
	}

	_ = w.Write([]string{}) // blank line

	// --- SECTION 2: Ingredients (market list) ----------------------
	_ = w.Write([]string{"# Ingredients (aggregated)"})
	_ = w.Write([]string{"Ingredient", "Unit", "Total Qty", "Unit Cost", "Total Cost"})

	// Sort by name for nicer output
	ingSlice := make([]*ingAgg, 0, len(ingredients))
	for _, v := range ingredients {
		ingSlice = append(ingSlice, v)
	}
	sort.Slice(ingSlice, func(i, j int) bool {
		return ingSlice[i].Name < ingSlice[j].Name
	})

	for _, ing := range ingSlice {
		_ = w.Write([]string{
			ing.Name,
			ing.Unit,
			fmt.Sprintf("%.3f", ing.TotalQty),
			fmt.Sprintf("%.2f", ing.CostPerUnit),
			fmt.Sprintf("%.2f", ing.TotalCost),
		})
	}

	_ = w.Write([]string{})

	// --- SECTION 3: Subrecipes (aggregated) ------------------------
	_ = w.Write([]string{"# Subrecipes (aggregated, for bulk prep)"})
	_ = w.Write([]string{"Subrecipe", "Unit", "Total Qty"})

	subSlice := make([]*subAgg, 0, len(subrecipes))
	for _, v := range subrecipes {
		subSlice = append(subSlice, v)
	}
	sort.Slice(subSlice, func(i, j int) bool {
		return subSlice[i].Name < subSlice[j].Name
	})

	for _, s := range subSlice {
		_ = w.Write([]string{
			s.Name,
			s.Unit,
			fmt.Sprintf("%.3f", s.TotalQty),
		})
	}

	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "error writing csv: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Forecast exported to %s\n", *outFile)
}
