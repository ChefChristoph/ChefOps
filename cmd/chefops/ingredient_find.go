package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ChefChristoph/chefops/internal"
)

func ingredientFind(args []string) {
	fs := flag.NewFlagSet("ingredient find", flag.ExitOnError)
	fs.Parse(args)

	if len(fs.Args()) == 0 {
		fmt.Println("usage: chefops ingredient find <search>")
		os.Exit(1)
	}

	search := "%" + fs.Args()[0] + "%"

	db, _ := internal.OpenDB()
	defer db.Close()

	const q = `
		SELECT id, name, unit, cost_per_unit
		FROM ingredients
		WHERE name LIKE ?
		ORDER BY name;
	`

	rows, err := db.Query(q, search)
	if err != nil {
		fmt.Fprintf(os.Stderr, "search error: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tUNIT\tCOST/UNIT")
	for rows.Next() {
		var id int
		var name, unit string
		var cost float64

		rows.Scan(&id, &name, &unit, &cost)
		fmt.Fprintf(w, "%d\t%s\t%s\t%.2f\n", id, name, unit, cost)
	}
	w.Flush()
}
