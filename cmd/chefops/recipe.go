package main

import (
	"flag"
	"fmt"
	"os"
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
// recipeList start
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

	// Clean, copy-safe output
	fmt.Printf("%-4s %-45s %-7s %s\n", "ID", "NAME", "YIELD", "UNIT")
	for rows.Next() {
		var id int
		var name, unit string
		var qty float64
		rows.Scan(&id, &name, &qty, &unit)

		// Left-align name but without trailing padding
		fmt.Printf("%-4d %-45s %-7.2f %s\n", id, name, qty, unit)
	}
}
// recipeList end
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

    raw := strings.Join(args, " ")
    db, _ := internal.OpenDB()
    defer db.Close()

    recipeID, recipeName, err := findRecipeByName(db, raw)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("recipe not found:", raw)
            os.Exit(1)
        }
        fmt.Println("error finding recipe:", err)
        os.Exit(1)
    }

    fmt.Printf("\nRecipe: %s\n", recipeName)
    fmt.Println("-----------------------------------")
    fmt.Println("| Type       | Name               | Qty     | Unit | Line cost |")
    fmt.Println("|------------|--------------------|---------|------|-----------|")

    rows, err := db.Query(`
        SELECT type, name, qty, unit, line_cost
        FROM recipe_raw_lines
        WHERE recipe_id = ?
        ORDER BY type, name;
    `, recipeID)
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
// func recipeCost start
func recipeCost(args []string) {
    if len(args) == 0 {
        fmt.Println("usage: chefops recipe cost NAME")
        os.Exit(1)
    }

    raw := strings.Join(args, " ")

    db, _ := internal.OpenDB()
    defer db.Close()

    recipeID, recipeName, err := findRecipeByName(db, raw)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("recipe not found:", raw)
            os.Exit(1)
        }
        fmt.Println("error finding recipe:", err)
        os.Exit(1)
    }

    _ = recipeID   // silence unused var; we'll use this later in forecasting
const q = `
        SELECT
            recipe_name,
            yield_qty,
            yield_unit,
            secondary_yield_qty,
            secondary_yield_unit,
            total_cost,
            cost_per_yield_unit,
            cost_per_secondary_unit
        FROM recipe_totals
        WHERE recipe_name = ?
    `

    var unit, secUnit string
    var yield, secYield, total, perYield float64
    var perSecYield sql.NullFloat64

    err = db.QueryRow(q, recipeName).Scan(
        &recipeName,
        &yield,
        &unit,
        &secYield,
        &secUnit,
        &total,
        &perYield,
        &perSecYield,
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "error calculating cost: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("\nCost Breakdown for: %s\n", recipeName)
    fmt.Println("-----------------------------------")
    fmt.Printf("Total Cost:          %.2f\n", total)

    fmt.Printf("Yield:               %.2f %s\n", yield, unit)
    fmt.Printf("Cost per %s:         %.4f\n", unit, perYield)

    if secUnit != "" && secYield > 0 && perSecYield.Valid {
        fmt.Printf("Secondary Yield:     %.2f %s\n", secYield, secUnit)
        fmt.Printf("Cost per %s:         %.4f\n", secUnit, perSecYield.Float64)
    }

    fmt.Println()
}
// func recipeCost end 


// helper func findRecipeByName

func findRecipeByName(db *sql.DB, input string) (int, string, error) {
    name := strings.TrimSpace(input)

    // 1) exact, case-sensitive
    var id int
    err := db.QueryRow(`SELECT id, name FROM recipes WHERE name = ?`, name).Scan(&id, &name)
    if err == nil {
        return id, name, nil
    }

    // 2) exact, case-insensitive
    err = db.QueryRow(`SELECT id, name FROM recipes WHERE LOWER(name) = LOWER(?)`, name).Scan(&id, &name)
    if err == nil {
        return id, name, nil
    }

    // 3) fuzzy: contains (case-insensitive)
    rows, err := db.Query(`
        SELECT id, name
        FROM recipes
        WHERE LOWER(name) LIKE '%' || LOWER(?) || '%'
        ORDER BY name
        LIMIT 10;
    `, name)
    if err != nil {
        return 0, "", err
    }
    defer rows.Close()

    var matches []struct {
        ID   int
        Name string
    }

    for rows.Next() {
        var mid int
        var mname string
        rows.Scan(&mid, &mname)
        matches = append(matches, struct {
            ID   int
            Name string
        }{mid, mname})
    }

    if len(matches) == 0 {
        return 0, "", sql.ErrNoRows
    }
    if len(matches) == 1 {
        return matches[0].ID, matches[0].Name, nil
    }

    fmt.Println("Multiple recipes match:")
    for i, m := range matches {
        fmt.Printf("  %d) %s\n", i+1, m.Name)
    }
    fmt.Print("Choose number or press Enter to cancel: ")

    var choiceStr string
    fmt.Scanln(&choiceStr)
    if choiceStr == "" {
        return 0, "", fmt.Errorf("cancelled")
    }

    choice, convErr := strconv.Atoi(choiceStr)
    if convErr != nil || choice < 1 || choice > len(matches) {
        return 0, "", fmt.Errorf("invalid choice")
    }

    return matches[choice-1].ID, matches[choice-1].Name, nil
}

// ------------------------------------------------------------
// RECIPE SCALE — flexible argument order (fully fixed)
// ------------------------------------------------------------

func recipeScale(args []string) {
    if len(args) == 0 {
        fmt.Println("usage: chefops recipe scale NAME --qty X --unit UNIT")
        os.Exit(1)
    }

    // 1) Extract recipe name first (all tokens until first --flag)
    recipeParts := []string{}
    flagStart := -1

    for i, a := range args {
        if strings.HasPrefix(a, "--") {
            flagStart = i
            break
        }
        recipeParts = append(recipeParts, a)
    }

    // Safety: if no flags detected → error
    if flagStart == -1 {
        fmt.Println("qty and unit required")
        os.Exit(1)
    }

    recipeNameInput := strings.Join(recipeParts, " ")

    // 2) Parse flags AFTER the recipe name
    fs := flag.NewFlagSet("recipe scale", flag.ExitOnError)
    qty := fs.Float64("qty", 0, "target yield quantity")
    unit := fs.String("unit", "", "target yield unit")
    fs.Parse(args[flagStart:])

    if *qty <= 0 || *unit == "" {
        fmt.Println("qty and unit required")
        os.Exit(1)
    }

    db, _ := internal.OpenDB()
    defer db.Close()

    recipeID, recipeName, err := findRecipeByName(db, recipeNameInput)
    if err != nil {
        fmt.Println("recipe not found:", recipeNameInput)
        os.Exit(1)
    }

    // Fetch base yield
    var baseQty float64
    var baseUnit string
    err = db.QueryRow(`
        SELECT yield_qty, yield_unit
        FROM recipes WHERE id = ?
    `, recipeID).Scan(&baseQty, &baseUnit)
    if err != nil {
        fmt.Println("error loading base yield:", err)
        os.Exit(1)
    }

    factor := *qty / baseQty

    // UI
    fmt.Printf("\n📐 Scale Recipe: %s\n", recipeName)
    fmt.Println("--------------------------------------")
    fmt.Printf("Original Yield: %.3f %s\n", baseQty, baseUnit)
    fmt.Printf("Target Yield:   %.3f %s\n", *qty, *unit)
    fmt.Printf("Scale Factor:   %.3f\n\n", factor)

    // Query expanded lines
    rows, err := db.Query(`
        SELECT item_type, ingredient_name, qty, unit
        FROM recipe_items_expanded_detail
        WHERE recipe_id = ?
        ORDER BY item_type, ingredient_name
    `, recipeID)

    if err != nil {
        fmt.Println("error loading recipe:", err)
        os.Exit(1)
    }
    defer rows.Close()

    fmt.Printf("| %-10s | %-20s | %-10s | %-6s |\n", "Type", "Name", "New Qty", "Unit")
    fmt.Println("|------------|----------------------|------------|--------|")

    for rows.Next() {
        var itemType, name, u string
        var q float64
        if err := rows.Scan(&itemType, &name, &q, &u); err != nil {
            fmt.Println("scan error:", err)
            os.Exit(1)
        }

        fmt.Printf("| %-10s | %-20s | %-10.3f | %-6s |\n",
            itemType, name, q*factor, u)
    }

    fmt.Println()
}
