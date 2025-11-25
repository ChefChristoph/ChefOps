package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/ChefChristoph/chefops/internal"
)

// ----------------------------------------------------------------------
// ADD SUBRECIPE
// ----------------------------------------------------------------------
func recipeAddSubrecipe(args []string) {
    fs := flag.NewFlagSet("recipe add-subrecipe", flag.ExitOnError)
    recipeName := fs.String("recipe", "", "parent recipe")
    subName := fs.String("sub", "", "subrecipe name")
    qty := fs.Float64("qty", 0, "quantity of subrecipe used")
    unit := fs.String("unit", "", "unit (optional; defaults to parent recipe unit)")
    fs.Parse(args)

    if *recipeName == "" || *subName == "" || *qty <= 0 {
        fmt.Println("usage: --recipe NAME --sub NAME --qty X [--unit unit]")
        os.Exit(1)
    }

    db, _ := internal.OpenDB()
    defer db.Close()

    // ---------------------------------------------------------
    // Look up parent recipe and default unit
    // ---------------------------------------------------------
    var recipeID int
    var parentUnit string

    err := db.QueryRow(`
        SELECT id, yield_unit
        FROM recipes
        WHERE name = ?
    `, *recipeName).Scan(&recipeID, &parentUnit)

    if err != nil {
        fmt.Println("recipe not found:", *recipeName)
        os.Exit(1)
    }

    // ---------------------------------------------------------
    // Subrecipe lookup
    // ---------------------------------------------------------
    var subID int
    err = db.QueryRow(`SELECT id FROM recipes WHERE name = ?`, *subName).Scan(&subID)
    if err != nil {
        fmt.Println("subrecipe not found:", *subName)
        os.Exit(1)
    }

    if recipeID == subID {
        fmt.Println("A recipe cannot reference itself.")
        os.Exit(1)
    }

    // ---------------------------------------------------------
    // Select effective unit (explicit > inherited)
    // ---------------------------------------------------------
    effectiveUnit := parentUnit
    if *unit != "" {
        effectiveUnit = *unit
    }

    // ---------------------------------------------------------
    // Insert
    // ---------------------------------------------------------
    _, err = db.Exec(`
        INSERT INTO recipe_subrecipes (recipe_id, subrecipe_id, qty, unit)
        VALUES (?, ?, ?, ?)
    `, recipeID, subID, *qty, effectiveUnit)

    if err != nil {
        fmt.Println("error adding subrecipe:", err)
        os.Exit(1)
    }

    fmt.Printf("Added subrecipe %s (%.3f %s) to %s\n",
        *subName, *qty, effectiveUnit, *recipeName)
}
// ----------------------------------------------------------------------
// REMOVE SUBRECIPE
// ----------------------------------------------------------------------
func recipeRemoveSubrecipe(args []string) {
	fs := flag.NewFlagSet("recipe remove-subrecipe", flag.ExitOnError)
	recipeName := fs.String("recipe", "", "parent recipe")
	subName := fs.String("sub", "", "subrecipe name")
	fs.Parse(args)

	if *recipeName == "" || *subName == "" {
		fmt.Println("usage: --recipe NAME --sub NAME")
		os.Exit(1)
	}

	db, _ := internal.OpenDB()
	defer db.Close()

	var recipeID int
	err := db.QueryRow(`SELECT id FROM recipes WHERE name = ?`, *recipeName).Scan(&recipeID)
	if err != nil {
		fmt.Println("recipe not found:", *recipeName)
		os.Exit(1)
	}

	var subID int
	err = db.QueryRow(`SELECT id FROM recipes WHERE name = ?`, *subName).Scan(&subID)
	if err != nil {
		fmt.Println("subrecipe not found:", *subName)
		os.Exit(1)
	}

	_, err = db.Exec(`
		DELETE FROM recipe_subrecipes
		WHERE recipe_id = ? AND subrecipe_id = ?
	`, recipeID, subID)

	if err != nil {
		fmt.Println("error removing subrecipe:", err)
		os.Exit(1)
	}

	fmt.Printf("Removed subrecipe %s from %s\n", *subName, *recipeName)
}

// ----------------------------------------------------------------------
// STRUCT FOR RETURNING SUBRECIPE LIST
// ----------------------------------------------------------------------
type SubrecipeEntry struct {
	Name string
	Qty  float64
	Unit string
}

// ----------------------------------------------------------------------
// LIST SUBRECIPES FOR A RECIPE
// (used by recipe show and export)
// ----------------------------------------------------------------------
func recipeListSubrecipes(recipeID int, db *sql.DB) ([]SubrecipeEntry, error) {

	rows, err := db.Query(`
		SELECT r.name, rs.qty, rs.unit
		FROM recipe_subrecipes rs
		JOIN recipes r ON r.id = rs.subrecipe_id
		WHERE rs.recipe_id = ?
		ORDER BY r.name
	`, recipeID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SubrecipeEntry

	for rows.Next() {
		var e SubrecipeEntry
		rows.Scan(&e.Name, &e.Qty, &e.Unit)
		list = append(list, e)
	}

	return list, nil
}
