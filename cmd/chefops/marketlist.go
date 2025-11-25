package main

import (
    "flag"
    "fmt"
    "os"
    "text/tabwriter"

    "github.com/ChefChristoph/chefops/internal"
)

func marketlist(args []string) {
    fs := flag.NewFlagSet("marketlist", flag.ExitOnError)
    fs.Parse(args)

    db, _ := internal.OpenDB()
    defer db.Close()

    const q = `
        SELECT
            ingredient_name,
            unit,
            cost_per_unit,
            total_qty,
            total_cost
        FROM market_list
        ORDER BY ingredient_name;
    `

    rows, err := db.Query(q)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error querying market list: %v\n", err)
        os.Exit(1)
    }
    defer rows.Close()

    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "INGREDIENT\tUNIT\tTOTAL QTY\tUNIT COST\tTOTAL COST")

    for rows.Next() {
        var name, unit string
        var qty, unitCost, totalCost float64

        if err := rows.Scan(&name, &unit, &unitCost, &qty, &totalCost); err != nil {
            fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
            os.Exit(1)
        }

        fmt.Fprintf(
            w, "%s\t%s\t%.3f\t%.2f\t%.2f\n",
            name, unit, qty, unitCost, totalCost,
        )
    }

    w.Flush()
}
