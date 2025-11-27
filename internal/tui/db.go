package tui

import (
	"database/sql"
	"fmt"
)

type RecipeSummary struct {
	ID        int
	Name      string
	TotalCost float64
}

type RecipeLine struct {
	Type string
	Name string
	Qty  float64
	Unit string
	Cost float64
}

type RecipeDetail struct {
	ID        int
	Name      string
	Lines     []RecipeLine
	TotalCost float64

	YieldQty  float64
	YieldUnit string
}

func LoadRecipes(db *sql.DB) ([]RecipeSummary, error) {
	rows, err := db.Query(`SELECT recipe_id, recipe_name, total_cost FROM recipe_totals ORDER BY recipe_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []RecipeSummary
	for rows.Next() {
		var r RecipeSummary
		rows.Scan(&r.ID, &r.Name, &r.TotalCost)
		list = append(list, r)
	}
	return list, nil
}

func LoadRecipeDetail(db *sql.DB, recipeID int) (*RecipeDetail, error) {
	var name, unit string
	var yield float64

	err := db.QueryRow(`
		SELECT name, yield_qty, yield_unit
		FROM recipes
		WHERE id = ?`, recipeID).Scan(&name, &yield, &unit)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT type, name, qty, unit, line_cost
		FROM recipe_raw_lines
		WHERE recipe_id = ?
		ORDER BY type, name
	`, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	detail := &RecipeDetail{
		ID:        recipeID,
		Name:      name,
		YieldQty:  yield,
		YieldUnit: unit,
	}

	for rows.Next() {
		var l RecipeLine
		rows.Scan(&l.Type, &l.Name, &l.Qty, &l.Unit, &l.Cost)
		detail.Lines = append(detail.Lines, l)
		detail.TotalCost += l.Cost
	}

	return detail, nil
}

func LoadRecipeNotes(db *sql.DB, recipeID int) (string, error) {
	var notes string
	query := "SELECT notes FROM recipes WHERE id = ?"
	err := db.QueryRow(query, recipeID).Scan(&notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("recipe with ID %d not found", recipeID)
		}
		return "", fmt.Errorf("failed to load recipe notes: %w", err)
	}
	return notes, nil
}
