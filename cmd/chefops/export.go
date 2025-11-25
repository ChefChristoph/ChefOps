package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ChefChristoph/chefops/internal"
)

///////////////////////////////////////////////////////////////////////////////
// EXPORT DISPATCHER
///////////////////////////////////////////////////////////////////////////////

func exportCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: chefops export <recipe|marketlist|full-report> [...]")
		os.Exit(1)
	}

	switch args[0] {
	case "recipe":
		exportRecipe(args[1:])
	case "marketlist":
		exportMarketlist(args[1:])
	case "full-report":
		exportFullReport(args[1:])
	default:
		fmt.Println("unknown export type:", args[0])
		os.Exit(1)
	}
}

///////////////////////////////////////////////////////////////////////////////
// FLAG PARSING (THE FIXED VERSION)
///////////////////////////////////////////////////////////////////////////////

type exportOptions struct {
	outfile string
	json    bool
}

func parseExportFlags(args []string) (exportOptions, []string) {
	opts := exportOptions{}
	positional := []string{}

	// manual flag parsing — bulletproof and simple
	for i := 0; i < len(args); i++ {
		a := args[i]

		// -o filename.md
		if a == "-o" && i+1 < len(args) {
			opts.outfile = args[i+1]
			i++
			continue
		}

		// --json
		if a == "--json" {
			opts.json = true
			continue
		}

		// positional argument
		positional = append(positional, a)
	}

	return opts, positional
}

///////////////////////////////////////////////////////////////////////////////
// EXPORT RECIPE
///////////////////////////////////////////////////////////////////////////////

func exportRecipe(args []string) {
	opts, positional := parseExportFlags(args)

	if len(positional) < 1 {
		fmt.Println("usage: chefops export recipe \"Recipe Name\" -o file.md")
		os.Exit(1)
	}

	recipeName := positional[0]

	db, _ := internal.OpenDB()
	defer db.Close()

	var yield, secYield float64
	var unit, secUnit string

	err := db.QueryRow(`
		SELECT yield_qty, yield_unit, secondary_yield_qty, secondary_yield_unit
		FROM recipes
		WHERE name = ?
	`, recipeName).Scan(&yield, &unit, &secYield, &secUnit)
	if err != nil {
		fmt.Println("recipe not found:", recipeName)
		os.Exit(1)
	}

	var totalCost, costPerYield float64
	var costPerSecondary sql.NullFloat64

	err = db.QueryRow(`
		SELECT total_cost, cost_per_yield_unit, cost_per_secondary_unit
		FROM recipe_totals
		WHERE recipe_name = ?
	`, recipeName).Scan(&totalCost, &costPerYield, &costPerSecondary)
	if err != nil {
		fmt.Println("error loading totals:", err)
		os.Exit(1)
	}

	rows, err := db.Query(`
		SELECT item_type, ingredient_name, qty, ingredient_unit, line_cost
		FROM recipe_items_expanded_detail_export
		WHERE recipe_name = ?
		ORDER BY item_type, ingredient_name
	`, recipeName)
	if err != nil {
		fmt.Println("error loading recipe items:", err)
		return
	}
	defer rows.Close()

	type line struct {
		Type     string  `json:"type"`
		Name     string  `json:"name"`
		Qty      float64 `json:"qty"`
		Unit     string  `json:"unit"`
		LineCost float64 `json:"line_cost"`
	}

	var lines []line
	for rows.Next() {
		var l line
		if err := rows.Scan(&l.Type, &l.Name, &l.Qty, &l.Unit, &l.LineCost); err != nil {
			fmt.Println("scan error:", err)
			continue
		}
		lines = append(lines, l)
	}

	// JSON EXPORT
	if opts.json {
		obj := map[string]interface{}{
			"recipe": recipeName,
			"yield": map[string]interface{}{
				"qty":  yield,
				"unit": unit,
			},
			"secondary_yield": map[string]interface{}{
				"qty":  secYield,
				"unit": secUnit,
			},
			"ingredients": lines,
			"totals": map[string]interface{}{
				"total_cost":      totalCost,
				"cost_per_unit":   costPerYield,
				"cost_per_second": costPerSecondary.Float64,
			},
		}

		data, _ := json.MarshalIndent(obj, "", "  ")
		writeOutput(opts.outfile, string(data))
		return
	}

	// MARKDOWN EXPORT
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", recipeName))
	sb.WriteString(fmt.Sprintf("**Yield:** %.2f %s\n\n", yield, unit))
	if secUnit != "" {
		sb.WriteString(fmt.Sprintf("**Secondary Yield:** %.2f %s\n\n", secYield, secUnit))
	}

	sb.WriteString("## Ingredients\n\n")
	sb.WriteString("| Type | Ingredient | Qty | Unit | Line Cost |\n")
	sb.WriteString("|------|------------|-----|------|-----------|\n")

	for _, l := range lines {
		sb.WriteString(fmt.Sprintf(
			"| %s | %s | %.3f | %s | %.2f |\n",
			l.Type, l.Name, l.Qty, l.Unit, l.LineCost,
		))
	}

	sb.WriteString("\n## Cost Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Cost:** %.2f\n", totalCost))
	sb.WriteString(fmt.Sprintf("- **Cost per %s:** %.4f\n", unit, costPerYield))
	if costPerSecondary.Valid {
		sb.WriteString(fmt.Sprintf("- **Cost per %s:** %.4f\n", secUnit, costPerSecondary.Float64))
	}

	writeOutput(opts.outfile, sb.String())
}

///////////////////////////////////////////////////////////////////////////////
// EXPORT MARKET LIST
///////////////////////////////////////////////////////////////////////////////

func exportMarketlist(args []string) {
	opts, _ := parseExportFlags(args)

	db, _ := internal.OpenDB()
	defer db.Close()

	rows, err := db.Query(`
		SELECT ingredient_name, total_qty, unit, cost_per_unit, total_cost
		FROM market_list
		ORDER BY ingredient_name
	`)
	if err != nil {
		fmt.Println("error loading market list:", err)
		return
	}
	defer rows.Close()

	type item struct {
		Name string  `json:"name"`
		Qty  float64 `json:"qty"`
		Unit string  `json:"unit"`
		Cost float64 `json:"cost"`
		Est  float64 `json:"estimated_cost"`
	}

	var list []item
	for rows.Next() {
		var it item
		rows.Scan(&it.Name, &it.Qty, &it.Unit, &it.Cost, &it.Est)
		list = append(list, it)
	}

	if opts.json {
		data, _ := json.MarshalIndent(list, "", "  ")
		writeOutput(opts.outfile, string(data))
		return
	}

	var sb strings.Builder

	sb.WriteString("# Market List\n\n")
	sb.WriteString("| Ingredient | Qty | Unit | Cost/Unit | Est Cost |\n")
	sb.WriteString("|-----------|-----|------|-----------|----------|\n")

	for _, it := range list {
		sb.WriteString(fmt.Sprintf(
			"| %s | %.3f | %s | %.2f | %.2f |\n",
			it.Name, it.Qty, it.Unit, it.Cost, it.Est,
		))
	}

	writeOutput(opts.outfile, sb.String())
}

///////////////////////////////////////////////////////////////////////////////
// EXPORT FULL REPORT
///////////////////////////////////////////////////////////////////////////////

func exportFullReport(args []string) {
	opts, _ := parseExportFlags(args)

	db, _ := internal.OpenDB()
	defer db.Close()

	rows, err := db.Query(`SELECT name FROM recipes ORDER BY name`)
	if err != nil {
		fmt.Println("error loading recipes:", err)
		return
	}
	defer rows.Close()

	type line struct {
		Type string
		Name string
		Qty  float64
		Unit string
		Cost float64
	}

	type block struct {
		Name   string      `json:"name"`
		Yield  interface{} `json:"yield"`
		Lines  interface{} `json:"lines"`
		Totals interface{} `json:"totals"`
	}

	var all []block
	var sb strings.Builder

	if !opts.json {
		sb.WriteString("# ChefOps Full Report\n\n")
	}

	for rows.Next() {
		var name string
		rows.Scan(&name)

		var y, sy float64
		var u, su string

		db.QueryRow(`
			SELECT yield_qty, yield_unit, secondary_yield_qty, secondary_yield_unit
			FROM recipes WHERE name = ?
		`, name).Scan(&y, &u, &sy, &su)

		var total, cpu float64
		var cps sql.NullFloat64

		db.QueryRow(`
			SELECT total_cost, cost_per_yield_unit, cost_per_secondary_unit
			FROM recipe_totals WHERE recipe_name = ?
		`, name).Scan(&total, &cpu, &cps)

		r2, _ := db.Query(`
			SELECT item_type, ingredient_name, qty, ingredient_unit, line_cost
			FROM recipe_items_expanded_detail_export
			WHERE recipe_name = ?
		`, name)

		var lines []line
		for r2.Next() {
			var l line
			r2.Scan(&l.Type, &l.Name, &l.Qty, &l.Unit, &l.Cost)
			lines = append(lines, l)
		}
		r2.Close()

		if opts.json {
			all = append(all, block{
				Name: name,
				Yield: map[string]interface{}{
					"qty":  y,
					"unit": u,
					"secondary": map[string]interface{}{
						"qty":  sy,
						"unit": su,
					},
				},
				Lines: lines,
				Totals: map[string]interface{}{
					"total":           total,
					"cost_per_unit":   cpu,
					"cost_per_second": cps.Float64,
				},
			})
		} else {
			sb.WriteString(fmt.Sprintf("## %s\n\n", name))
			sb.WriteString(fmt.Sprintf("Yield: %.2f %s\n", y, u))
			if su != "" {
				sb.WriteString(fmt.Sprintf("Secondary Yield: %.2f %s\n", sy, su))
			}
			sb.WriteString("\n")
		}
	}

	if opts.json {
		data, _ := json.MarshalIndent(all, "", "  ")
		writeOutput(opts.outfile, string(data))
		return
	}

	writeOutput(opts.outfile, sb.String())
}

///////////////////////////////////////////////////////////////////////////////
// WRITE OUTPUT
///////////////////////////////////////////////////////////////////////////////

func writeOutput(outfile, content string) {
	if outfile == "" {
		fmt.Println(content)
		return
	}

	err := os.WriteFile(outfile, []byte(content), 0644)
	if err != nil {
		fmt.Println("error writing file:", err)
		return
	}

	fmt.Println("Saved:", outfile)
}
