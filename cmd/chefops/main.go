package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ChefChristoph/chefops/internal"
)
func usage() {
	    fmt.Println("ChefOps CLI")
    fmt.Println("")
    fmt.Println("Usage:")
    fmt.Println("  chefops ingredient add        --name NAME --unit UNIT --cost COST")
    fmt.Println("  chefops ingredient list")
    fmt.Println("")
    fmt.Println("  chefops recipe new            --name NAME --yield QTY --unit UNIT [--syield QTY --sunit UNIT]")
    fmt.Println("  chefops recipe list")
    fmt.Println("  chefops recipe show           \"RECIPE NAME\"")
    fmt.Println("  chefops recipe cost           \"RECIPE NAME\"")
    fmt.Println("  chefops recipe add-item       --recipe NAME --ingredient NAME --qty QTY")
    fmt.Println("  chefops recipe add-subrecipe  --recipe NAME --sub NAME --qty QTY --unit UNIT")
    fmt.Println("  chefops recipe scale          \"RECIPE NAME\" --qty QTY --unit UNIT")
    fmt.Println("")
    fmt.Println("  chefops forecast              \"DISH NAME\" --portions N")
    fmt.Println("")
    fmt.Println("  chefops marketlist")
    fmt.Println("")
    fmt.Println("Examples:")
    fmt.Println("  chefops recipe show \"BULK Batter\"")
    fmt.Println("  chefops recipe cost \"DISH Turbo Hammour Popcorn\"")
    fmt.Println("  chefops recipe scale \"BULK Batter\" --qty 8 --unit kg")
    fmt.Println("  chefops forecast \"DISH Pole Position Burger\" --portions 120")

	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {

	// -------------------
	// INGREDIENT COMMANDS
	// -------------------
	case "ingredient":
		if len(os.Args) < 3 {
			usage()
		}
		switch os.Args[2] {
		case "add":
			ingredientAdd(os.Args[3:])
		case "list":
			ingredientList(os.Args[3:])
		case "find":
			ingredientFind(os.Args[3:])
		case "convert":
		 ingredientConversionCommand(os.Args[3:])
		default:
			usage()
		}

	// --------------
	// RECIPE COMMANDS
	// --------------
	case "recipe":
		if len(os.Args) < 3 {
			usage()
		}
		switch os.Args[2] {
		case "new":
			recipeNew(os.Args[3:])
		case "list":
			recipeList(os.Args[3:])
		case "add-item":
			recipeAddItem(os.Args[3:])
		case "show":
			recipeShow(os.Args[3:])
		case "cost":
			recipeCost(os.Args[3:])
		case "remove-item":
      recipeRemoveItem(os.Args[3:])
		case "add-subrecipe":
    recipeAddSubrecipe(os.Args[3:])
case "remove-subrecipe":
    recipeRemoveSubrecipe(os.Args[3:])
		case "scale":
    recipeScale(os.Args[3:])
		default:
			usage()
		}
  // FORECAST COMMANDS
	case "forecast":
    forecastCommand(os.Args[2:])
	// -------------------------
  // MARKETLIST COMMANDS
	// -------------------------
	case "marketlist":
	marketlist(os.Args[2:])

  // -------------------------
  // export COMMANDS
	// -------------------------
	case "export":
	 exportCommand(os.Args[2:])

	// -------------------------
	// UNKNOWN TOP-LEVEL COMMAND
	// -------------------------
	default:
		usage()
	}
}
func openDBOrExit() *sql.DB {
	db, err := internal.OpenDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening DB: %v\n", err)
		os.Exit(1)
	}
	return db
}

func ingredientAdd(args []string) {
	fs := flag.NewFlagSet("ingredient add", flag.ExitOnError)
	name := fs.String("name", "", "ingredient name")
	unit := fs.String("unit", "", "unit, e.g. kg, l, piece")
	cost := fs.Float64("cost", 0, "cost per unit")
	fs.Parse(args)

	if *name == "" || *unit == "" || *cost <= 0 {
		fmt.Fprintln(os.Stderr, "name, unit and positive cost are required")
		fs.Usage()
		os.Exit(1)
	}

	db := openDBOrExit()
	defer db.Close()

	const q = `
		INSERT INTO ingredients (name, unit, cost_per_unit)
		VALUES (?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
		    unit = excluded.unit,
		    cost_per_unit = excluded.cost_per_unit;
	`
	_, err := db.Exec(q, *name, *unit, *cost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error inserting ingredient: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ingredient saved: %s (%s @ %.2f)\n", *name, *unit, *cost)
}

func ingredientList(args []string) {
	fs := flag.NewFlagSet("ingredient list", flag.ExitOnError)
	fs.Parse(args)

	db := openDBOrExit()
	defer db.Close()

	rows, err := db.Query(`SELECT id, name, unit, cost_per_unit FROM ingredients ORDER BY name;`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error querying ingredients: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tUNIT\tCOST/UNIT")
	for rows.Next() {
		var (
			id   int
			name string
			unit string
			cost float64
		)
		if err := rows.Scan(&id, &name, &unit, &cost); err != nil {
			fmt.Fprintf(os.Stderr, "error scanning row: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%.4f\n", id, name, unit, cost)
	}
	w.Flush()
}
