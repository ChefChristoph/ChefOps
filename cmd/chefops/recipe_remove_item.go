package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/ChefChristoph/chefops/internal"
)

func recipeRemoveItem(args []string) {
	fs := flag.NewFlagSet("recipe remove-item", flag.ExitOnError)
	recipeName := fs.String("recipe", "", "recipe name")
	ingredientName := fs.String("ingredient", "", "ingredient name (fuzzy)")
	fs.Parse(args)

	if *recipeName == "" || *ingredientName == "" {
		fmt.Println("usage: chefops recipe remove-item --recipe NAME --ingredient NAME")
		os.Exit(1)
	}

	db, _ := internal.OpenDB()
	defer db.Close()

	// --- Find recipe ---
	var recipeID int
	err := db.QueryRow(`SELECT id FROM recipes WHERE name = ?`, *recipeName).Scan(&recipeID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "recipe not found: %s\n", *recipeName)
		os.Exit(1)
	}

	// --- Find matching items inside recipe ---
	const q = `
		SELECT ri.id, i.name, ri.qty, i.unit
		FROM recipe_items ri
		JOIN ingredients i ON i.id = ri.ingredient_id
		WHERE ri.recipe_id = ?
		  AND LOWER(i.name) LIKE LOWER(?)
		ORDER BY i.name;
	`
	rows, err := db.Query(q, recipeID, "%"+*ingredientName+"%")
	if err != nil {
		fmt.Fprintf(os.Stderr, "query error: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	type item struct {
		ID    int
		Name  string
		Qty   float64
		Unit  string
	}
	var matches []item

	for rows.Next() {
		var it item
		rows.Scan(&it.ID, &it.Name, &it.Qty, &it.Unit)
		matches = append(matches, it)
	}

	// Nothing found?
	if len(matches) == 0 {
		fmt.Println("No matching items found in this recipe.")
		os.Exit(1)
	}

	// One match → confirm removal
	if len(matches) == 1 {
		fmt.Printf("Remove %.3f %s of %s from %s? (y/N): ",
			matches[0].Qty, matches[0].Unit, matches[0].Name, *recipeName)

		var choice string
		fmt.Scanln(&choice)

		if choice != "y" && choice != "Y" {
			fmt.Println("Cancelled.")
			os.Exit(0)
		}

		db.Exec(`DELETE FROM recipe_items WHERE id = ?`, matches[0].ID)
		fmt.Println("Removed.")
		return
	}

	// Multiple matches → interactive picker
	fmt.Println("Multiple items match that ingredient name:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "INDEX\tINGREDIENT\tQTY\tUNIT")
	for i, m := range matches {
		fmt.Fprintf(w, "%d\t%s\t%.3f\t%s\n", i+1, m.Name, m.Qty, m.Unit)
	}
	w.Flush()

	fmt.Printf("\nChoose number (1-%d) or Enter to cancel: ", len(matches))

	var pickStr string
	fmt.Scanln(&pickStr)
	if pickStr == "" {
		fmt.Println("Cancelled.")
		os.Exit(0)
	}

	pick, err := strconv.Atoi(pickStr)
	if err != nil || pick < 1 || pick > len(matches) {
		fmt.Println("Invalid selection.")
		os.Exit(1)
	}

	// Confirm delete
	chosen := matches[pick-1]
	fmt.Printf("Remove %.3f %s of %s from %s? (y/N): ",
		chosen.Qty, chosen.Unit, chosen.Name, *recipeName)

	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Cancelled.")
		return
	}

	db.Exec(`DELETE FROM recipe_items WHERE id = ?`, chosen.ID)
	fmt.Println("Removed.")
}
