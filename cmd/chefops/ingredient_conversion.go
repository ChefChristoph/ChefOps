package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/ChefChristoph/chefops/internal"
)

func ingredientConversionCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: chefops ingredient convert <add|list> ...")
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		ingredientConversionAdd(args[1:])
	case "list":
		ingredientConversionList(args[1:])
	default:
		fmt.Println("unknown conversion subcommand:", args[0])
		os.Exit(1)
	}
}

func ingredientConversionAdd(args []string) {
	fs := flag.NewFlagSet("ingredient convert add", flag.ExitOnError)
	ingredientName := fs.String("ingredient", "", "ingredient name")
	fromStr := fs.String("from", "", "from quantity + unit (e.g. 1kg)")
	toStr := fs.String("to", "", "to quantity + unit (e.g. 10piece)")
	fs.Parse(args)

	if *ingredientName == "" || *fromStr == "" || *toStr == "" {
		fmt.Println("usage: --ingredient NAME --from 1kg --to 10piece")
		os.Exit(1)
	}

	db, _ := internal.OpenDB()
	defer db.Close()

	// Get ingredient ID
	var ingID int
	err := db.QueryRow(`SELECT id FROM ingredients WHERE name = ?`, *ingredientName).Scan(&ingID)
	if err != nil {
		fmt.Println("ingredient not found:", *ingredientName)
		os.Exit(1)
	}

	parse := func(s string) (qty float64, unit string) {
		// Example input: 1kg, 10piece, 100g_breadcrumbs
		i := 0
		for ; i < len(s); i++ {
			if (s[i] < '0' || s[i] > '9') && s[i] != '.' {
				break
			}
		}
		qtyStr := s[:i]
		unit = s[i:]
		qty, _ = strconv.ParseFloat(qtyStr, 64)
		return
	}

	fromQty, fromUnit := parse(*fromStr)
	toQty, toUnit := parse(*toStr)

	_, err = db.Exec(`
		INSERT INTO ingredient_conversions 
		(ingredient_id, from_unit, from_qty, to_unit, to_qty)
		VALUES (?, ?, ?, ?, ?)
	`, ingID, fromUnit, fromQty, toUnit, toQty)

	if err != nil {
		fmt.Println("error adding conversion:", err)
		os.Exit(1)
	}

	fmt.Printf("Added conversion: %s → %s\n", *fromStr, *toStr)
}

func ingredientConversionList(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: chefops ingredient convert list \"Ingredient\"")
		os.Exit(1)
	}
	name := args[0]

	db, _ := internal.OpenDB()
	defer db.Close()

	var ingID int
	err := db.QueryRow(`SELECT id FROM ingredients WHERE name = ?`, name).Scan(&ingID)
	if err != nil {
		fmt.Println("ingredient not found:", name)
		os.Exit(1)
	}

	rows, _ := db.Query(`
		SELECT from_qty, from_unit, to_qty, to_unit
		FROM ingredient_conversions
		WHERE ingredient_id = ?
	`, ingID)
	defer rows.Close()

	fmt.Printf("\nConversions for %s:\n", name)
	for rows.Next() {
		var fq, tq float64
		var fu, tu string
		rows.Scan(&fq, &fu, &tq, &tu)
		fmt.Printf("  %.3f %s → %.3f %s\n", fq, fu, tq, tu)
	}
}
