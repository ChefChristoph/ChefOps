package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
  "strconv"
	"strings"
	"database/sql"
	"github.com/ChefChristoph/chefops/internal"
)

func recipeNew(args []string) {
	fs := flag.NewFlagSet("recipe new", flag.ExitOnError)
	name := fs.String("name", "", "recipe name")
	yieldQty := fs.Float64("yield", 0, "primary yield quantity (e.g. 10)")
	yieldUnit := fs.String("unit", "", "primary yield unit (e.g. kg, l, portion)")
	secYieldQty := fs.Float64("syield", 0, "secondary yield quantity (e.g. 40)")
	secYieldUnit := fs.String("sunit", "", "secondary yield unit (e.g. piece)")
	fs.Parse(args)

	if *name == "" {
		fmt.Fprintln(os.Stderr, "recipe name required")
		fs.Usage()
		os.Exit(1)
	}

	db, _ := internal.OpenDB()
	defer db.Close()

	const q = `
	INSERT INTO recipes (name, yield_qty, yield_unit, secondary_yield_qty, secondary_yield_unit)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(name) DO UPDATE SET
		yield_qty = excluded.yield_qty,
		yield_unit = excluded.yield_unit,
		secondary_yield_qty = excluded.secondary_yield_qty,
		secondary_yield_unit = excluded.secondary_yield_unit;
	`

	_, err := db.Exec(q, *name, *yieldQty, *yieldUnit, *secYieldQty, *secYieldUnit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating recipe: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Recipe saved: %s (yield %.2f %s", *name, *yieldQty, *yieldUnit)
	if *secYieldQty > 0 && *secYieldUnit != "" {
		fmt.Printf(", secondary yield %.2f %s", *secYieldQty, *secYieldUnit)
	}
	fmt.Println(")")
}

func recipeList(args []string) {
	fs := flag.NewFlagSet("recipe list", flag.ExitOnError)
	fs.Parse(args)

	db, _ := internal.OpenDB()
	defer db.Close()

	rows, err := db.Query(`SELECT id, name, yield_qty, yield_unit FROM recipes ORDER BY name;`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing recipes: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tYIELD\tUNIT")
	for rows.Next() {
		var id int
		var name, unit string
		var qty float64
		rows.Scan(&id, &name, &qty, &unit)
		fmt.Fprintf(w, "%d\t%s\t%.2f\t%s\n", id, name, qty, unit)
	}
	w.Flush()
}

func recipeAddItem(args []string) {
	fs := flag.NewFlagSet("recipe add-item", flag.ExitOnError)
	recipeName := fs.String("recipe", "", "recipe name")
	ingredientName := fs.String("ingredient", "", "ingredient name")
	qty := fs.Float64("qty", 0, "quantity")
	fs.Parse(args)

	if *recipeName == "" || *ingredientName == "" || *qty <= 0 {
		fmt.Fprintln(os.Stderr, "recipe, ingredient and positive qty are required")
		fs.Usage()
		os.Exit(1)
	}

	db, _ := internal.OpenDB()
	defer db.Close()

	// ---------------------------
	// Find recipe
	// ---------------------------
	var recipeID int
	err := db.QueryRow(`SELECT id FROM recipes WHERE name = ?`, *recipeName).Scan(&recipeID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "recipe not found: %s\n", *recipeName)
		os.Exit(1)
	}

	// We track the final ingredient name used for printing:
	actualIngredientName := *ingredientName

	// ---------------------------
	// Fuzzy search for ingredient
	// ---------------------------
	var ingredientID int
	err = db.QueryRow(`SELECT id FROM ingredients WHERE name = ?`, *ingredientName).Scan(&ingredientID)

	if err != nil {
		// Ingredient not found → fuzzy match
		input := *ingredientName
		fmt.Fprintf(os.Stderr, "ingredient not found: %s\n\n", input)
		fmt.Println("🔍 Searching for similar ingredients...")

		type suggestion struct {
			ID   int
			Name string
		}
		var suggestions []suggestion

		// Pass 1: contains pattern
		const q1 = `
			SELECT id, name FROM ingredients
			WHERE LOWER(name) LIKE LOWER(?)
			ORDER BY name LIMIT 10;
		`
		rows1, _ := db.Query(q1, "%"+input+"%")
		for rows1.Next() {
			var id int
			var name string
			rows1.Scan(&id, &name)
			suggestions = append(suggestions, suggestion{ID: id, Name: name})
		}
		rows1.Close()

		// Pass 2: prefix match (2 chars)
		if len(suggestions) == 0 && len(input) >= 2 {
			const q2 = `
				SELECT id, name FROM ingredients
				WHERE SUBSTR(LOWER(name),1,2) = SUBSTR(LOWER(?),1,2)
				ORDER BY name LIMIT 10;
			`
			rows2, _ := db.Query(q2, input)
			for rows2.Next() {
				var id int
				var name string
				rows2.Scan(&id, &name)
				suggestions = append(suggestions, suggestion{ID: id, Name: name})
			}
			rows2.Close()
		}

		// Pass 3: first word prefix match (fallback)
		if len(suggestions) == 0 {
			parts := strings.Split(strings.ToLower(input), " ")
			base := parts[0]

			const q3 = `
				SELECT id, name FROM ingredients
				WHERE LOWER(name) LIKE ?
				ORDER BY name LIMIT 10;
			`
			rows3, _ := db.Query(q3, base+"%")
			for rows3.Next() {
				var id int
				var name string
				rows3.Scan(&id, &name)
				suggestions = append(suggestions, suggestion{ID: id, Name: name})
			}
			rows3.Close()
		}

		if len(suggestions) == 0 {
			fmt.Println("(no similar ingredients found)")
			os.Exit(1)
		}

		// Auto-select if exactly one match
		if len(suggestions) == 1 {
			fmt.Printf("Using: %s\n", suggestions[0].Name)
			ingredientID = suggestions[0].ID
			actualIngredientName = suggestions[0].Name
		} else {
			// Interactive selection
			fmt.Println("Did you mean:")
			for i, s := range suggestions {
				fmt.Printf("  %d) %s\n", i+1, s.Name)
			}

			fmt.Printf("\nChoose number (1-%d) or press Enter to cancel: ", len(suggestions))

			var choiceStr string
			fmt.Scanln(&choiceStr)

			if choiceStr == "" {
				fmt.Println("Cancelled.")
				os.Exit(0)
			}

			choice, convErr := strconv.Atoi(choiceStr)
			if convErr != nil || choice < 1 || choice > len(suggestions) {
				fmt.Println("Invalid choice.")
				os.Exit(1)
			}

			ingredientID = suggestions[choice-1].ID
			actualIngredientName = suggestions[choice-1].Name
			fmt.Printf("Using: %s\n", actualIngredientName)
		}
	}

	// ---------------------------
	// Check for duplicate item
	// ---------------------------
	var existingQty float64
	err = db.QueryRow(`
		SELECT qty FROM recipe_items
		WHERE recipe_id = ? AND ingredient_id = ?
	`, recipeID, ingredientID).Scan(&existingQty)

	if err == nil {
		// Ingredient exists → choose add / replace
		fmt.Printf("\nIngredient '%s' already exists in recipe '%s'.\n", actualIngredientName, *recipeName)
		fmt.Printf("Current amount: %.3f\n\n", existingQty)

		fmt.Println("Choose:")
		fmt.Println("  1) Add to existing amount (+=)")
		fmt.Println("  2) Replace existing amount (=)")
		fmt.Println("  3) Cancel")

		fmt.Print("\nEnter choice: ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			newQty := existingQty + *qty
			_, err := db.Exec(`
				UPDATE recipe_items
				SET qty = ?
				WHERE recipe_id = ? AND ingredient_id = ?
			`, newQty, recipeID, ingredientID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "update error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Updated: %.3f → %.3f (added)\n", existingQty, newQty)
			return

		case "2":
			_, err := db.Exec(`
				UPDATE recipe_items
				SET qty = ?
				WHERE recipe_id = ? AND ingredient_id = ?
			`, *qty, recipeID, ingredientID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "update error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Updated: %.3f → %.3f (replaced)\n", existingQty, *qty)
			return

		case "3", "":
			fmt.Println("Cancelled.")
			return

		default:
			fmt.Println("Invalid choice. Cancelled.")
			return
		}
	}

	// ---------------------------
	// Insert new item (no duplicate)
	// ---------------------------
	_, err = db.Exec(`
		INSERT INTO recipe_items (recipe_id, ingredient_id, qty)
		VALUES (?, ?, ?)
	`, recipeID, ingredientID, *qty)

	if err != nil {
		fmt.Fprintf(os.Stderr, "insert error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added %.3f × %s to %s\n", *qty, actualIngredientName, *recipeName)
}
// func recipeShow start
func recipeShow(args []string) {
    if len(args) == 0 {
        fmt.Println("usage: chefops recipe show NAME")
        os.Exit(1)
    }

    name := args[0]

    db, _ := internal.OpenDB()
    defer db.Close()

    // Get recipe ID (optional but used earlier)
    var recipeID int
    err := db.QueryRow(`SELECT id FROM recipes WHERE name = ?`, name).Scan(&recipeID)
    if err != nil {
        fmt.Println("recipe not found:", name)
        os.Exit(1)
    }

    fmt.Printf("\nRecipe: %s\n", name)
    fmt.Println("-----------------------------------")
    fmt.Println("| Type       | Name               | Qty     | Unit | Line cost |")
    fmt.Println("|------------|--------------------|---------|------|-----------|")

    rows, err := db.Query(`
        SELECT item_type, ingredient_name, qty, ingredient_unit, line_cost
        FROM recipe_items_expanded
        WHERE recipe_name = ?
        ORDER BY item_type, ingredient_name
    `, name)

    if err != nil {
        fmt.Println("error loading recipe:", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var itemType, iname, unit string
        var qty, lineCost float64

        if err := rows.Scan(&itemType, &iname, &qty, &unit, &lineCost); err != nil {
            fmt.Println("scan error:", err)
            return
        }

        fmt.Printf(
            "| %-10s | %-18s | %-7.3f | %-4s | %-9.2f |\n",
            itemType, iname, qty, unit, lineCost,
        )
    }

    fmt.Println()
}

// func recipeShow end
func recipeCost(args []string) {
    if len(args) == 0 {
        fmt.Println("usage: chefops recipe cost NAME")
        os.Exit(1)
    }

    name := args[0]

    db, _ := internal.OpenDB()
    defer db.Close()

    const q = `
    SELECT recipe_name, yield_qty, yield_unit,
           secondary_yield_qty, secondary_yield_unit,
           total_cost, cost_per_yield_unit, cost_per_secondary_unit
    FROM recipe_totals
    WHERE recipe_name = ?
    `

    var recipeName, unit, secUnit string
    var yield, secYield, total, perUnit float64
    var perSecUnit sql.NullFloat64

    err := db.QueryRow(q, name).Scan(
        &recipeName,
        &yield, &unit,
        &secYield, &secUnit,
        &total, &perUnit,
        &perSecUnit,
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "error calculating cost: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("\nCost Breakdown for: %s\n", recipeName)
    fmt.Println("-----------------------------------")
    fmt.Printf("Total Cost:          %.2f\n", total)

    if yield > 0 && unit != "" {
        fmt.Printf("Yield:               %.2f %s\n", yield, unit)
        fmt.Printf("Cost per %s:         %.4f\n", unit, perUnit)
    }

    if secYield > 0 && secUnit != "" {
        fmt.Printf("Secondary Yield:     %.2f %s\n", secYield, secUnit)
        if perSecUnit.Valid {
            fmt.Printf("Cost per %s:         %.4f\n", secUnit, perSecUnit.Float64)
        }
    }

    fmt.Println()
}

