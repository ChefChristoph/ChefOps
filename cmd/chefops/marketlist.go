package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ChefChristoph/chefops/internal"
)

type MarketItem struct {
	IngredientID int     `json:"ingredient_id"`
	Name         string  `json:"ingredient_name"`
	Unit         string  `json:"unit"`
	TotalQty     float64 `json:"total_qty_needed"`
	CostPerUnit  float64 `json:"cost_per_unit"`
	EstCost      float64 `json:"estimated_cost"`
}

func marketList(args []string) {
	fs := flag.NewFlagSet("marketlist", flag.ExitOnError)
	markdown := fs.Bool("markdown", false, "output as markdown table")
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	db, _ := internal.OpenDB()
	defer db.Close()

	const q = `
	SELECT ingredient_id, ingredient_name, unit, total_qty_needed, cost_per_unit, estimated_cost
	FROM market_list
	ORDER BY ingredient_name;
	`

	rows, err := db.Query(q)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error querying market list: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var items []MarketItem

	for rows.Next() {
		var m MarketItem
		if err := rows.Scan(&m.IngredientID, &m.Name, &m.Unit, &m.TotalQty, &m.CostPerUnit, &m.EstCost); err != nil {
			fmt.Fprintf(os.Stderr, "error scanning row: %v\n", err)
			os.Exit(1)
		}
		items = append(items, m)
	}

	// JSON OUTPUT
	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(items)
		return
	}

	// MARKDOWN OUTPUT
	if *markdown {
		fmt.Println("| Ingredient | Total Qty | Unit | Cost/Unit | Est. Cost |")
		fmt.Println("|------------|-----------|------|-----------|-----------|")
		for _, m := range items {
			fmt.Printf("| %s | %.3f | %s | %.2f | %.2f |\n",
				m.Name, m.TotalQty, m.Unit, m.CostPerUnit, m.EstCost)
		}
		return
	}

	// DEFAULT: Table in terminal
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "INGREDIENT\tTOTAL QTY\tUNIT\tCOST/U\tEST. COST\n")
	for _, m := range items {
		fmt.Fprintf(w, "%s\t%.3f\t%s\t%.2f\t%.2f\n",
			m.Name, m.TotalQty, m.Unit, m.CostPerUnit, m.EstCost)
	}
	w.Flush()
}
