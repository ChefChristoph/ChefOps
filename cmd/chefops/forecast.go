package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ChefChristoph/chefops/internal"
)

func forecastCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("usage: chefops forecast RECIPE --portions X")
		os.Exit(1)
	}

	// ---------------------
	// Parse recipe name first
	// ---------------------
	recipeParts := []string{}
	flagStart := 0

	for i, a := range args {
		if strings.HasPrefix(a, "--") {
			flagStart = i
			break
		}
		recipeParts = append(recipeParts, a)
	}

	if len(recipeParts) == 0 {
		fmt.Println("usage: chefops forecast RECIPE --portions X")
		os.Exit(1)
	}

	recipeNameInput := strings.Join(recipeParts, " ")

	// ---------------------
	// Parse flags
	// ---------------------
	fs := flag.NewFlagSet("forecast", flag.ExitOnError)
	portions := fs.Float64("portions", 0, "number of portions")
	fs.Parse(args[flagStart:])

	if *portions <= 0 {
		fmt.Println("--portions required and must be > 0")
		os.Exit(1)
	}

	// ---------------------
	// DB
	// ---------------------
	db, _ := internal.OpenDB()
	defer db.Close()

	recipeID, recipeName, err := findRecipeByName(db, recipeNameInput)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("recipe not found:", recipeNameInput)
			os.Exit(1)
		}
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// ---------------------
	// Header
	// ---------------------
	fmt.Printf("\n📊 Forecast: %s\n", recipeName)
	fmt.Println("-----------------------------------------")
	fmt.Printf("Portions Required: %.2f\n\n", *portions)

	// ---------------------
	// Part 1 — Subrecipes
	// ---------------------
	rows, err := db.Query(`
        SELECT item_type, ingredient_name, qty, unit
        FROM recipe_items_expanded_detail
        WHERE recipe_id = ?
        ORDER BY item_type, ingredient_name
    `, recipeID)

	if err != nil {
		fmt.Println("error loading expanded recipe:", err)
		os.Exit(1)
	}
	defer rows.Close()

	type Line struct {
		Type string
		Name string
		Qty  float64
		Unit string
	}

	var lines []Line

	for rows.Next() {
		var t, n, u string
		var q float64

		if err := rows.Scan(&t, &n, &q, &u); err != nil {
			fmt.Println("scan error:", err)
			os.Exit(1)
		}

		lines = append(lines, Line{t, n, q, u})
	}

	// ---------------------
	// Output Subrecipes
	// ---------------------
	fmt.Println("🔧 Subrecipes required:")
	fmt.Println("-----------------------------------------")
	fmt.Printf("| %-25s | %-12s | %-6s |\n", "Subrecipe", "Qty Needed", "Unit")
	fmt.Println("|---------------------------|--------------|--------|")

	for _, l := range lines {
		if l.Type == "subrecipe" {
			fmt.Printf("| %-25s | %-12.3f | %-6s |\n",
				l.Name, l.Qty**portions, l.Unit)
		}
	}

	fmt.Println()

	// ---------------------
	// Output Ingredients
	// ---------------------
	fmt.Println("🥕 Ingredients required:")
	fmt.Println("-----------------------------------------")
	fmt.Printf("| %-25s | %-12s | %-6s |\n", "Ingredient", "Qty Needed", "Unit")
	fmt.Println("|---------------------------|--------------|--------|")

	for _, l := range lines {
		if l.Type == "ingredient" {
			fmt.Printf("| %-25s | %-12.3f | %-6s |\n",
				l.Name, l.Qty**portions, l.Unit)
		}
	}

	fmt.Println()
}

