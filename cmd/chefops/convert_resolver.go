package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ConversionStep struct {
	FromUnit string
	ToUnit   string
	Factor   float64 // multiply from_qty to get to_qty ratio
}

// ---------------------------------------------------------------------
// Resolve cost for an ingredient with (qty, unit)
// ---------------------------------------------------------------------
func resolveIngredientCost(ingredientID int, qty float64, unit string, db *sql.DB) (float64, string, error) {

	// Get base ingredient cost
	var baseUnit string
	var costPerUnit float64

	err := db.QueryRow(`
		SELECT unit, cost_per_unit
		FROM ingredients
		WHERE id = ?
	`, ingredientID).Scan(&baseUnit, &costPerUnit)

	if err != nil {
		return 0, "", fmt.Errorf("ingredient not found")
	}

	// If the unit matches the base unit, easy:
	// cost = qty * cost/kg
	if unit == baseUnit {
		return qty * costPerUnit, baseUnit, nil
	}

	// Try to find conversion path
	convertedQty, finalUnit, err := resolveConversionChain(ingredientID, qty, unit, baseUnit, db, map[string]bool{})
	if err != nil {
		return 0, "", err
	}

	// Now cost = convertedQty * cost_per_unit
	return convertedQty * costPerUnit, finalUnit, nil
}

// ---------------------------------------------------------------------
// Resolve conversion chain: unit → baseUnit, recursively
// ---------------------------------------------------------------------
func resolveConversionChain(ingID int, qty float64, fromUnit, targetUnit string, db *sql.DB, visited map[string]bool) (float64, string, error) {

	// Loop protection
	if visited[fromUnit] {
		return 0, "", errors.New("conversion loop detected")
	}
	visited[fromUnit] = true

	// If we reached target unit → done
	if fromUnit == targetUnit {
		return qty, targetUnit, nil
	}

	// Look for direct conversions
	rows, err := db.Query(`
		SELECT from_qty, from_unit, to_qty, to_unit
		FROM ingredient_conversions
		WHERE ingredient_id = ?
		  AND from_unit = ?
	`, ingID, fromUnit)
	if err != nil {
		return 0, "", err
	}
	defer rows.Close()

	type conv struct {
		fQty   float64
		fUnit  string
		tQty   float64
		tUnit  string
	}
	var convs []conv

	for rows.Next() {
		var c conv
		rows.Scan(&c.fQty, &c.fUnit, &c.tQty, &c.tUnit)
		convs = append(convs, c)
	}

	// No conversions from this unit
	if len(convs) == 0 {
		return 0, "", fmt.Errorf("no conversion from unit '%s' for ingredient", fromUnit)
	}

	// Try each conversion path
	for _, c := range convs {

		// ratio = to_qty / from_qty
		ratio := c.tQty / c.fQty
		nextQty := qty * ratio

		// Recursively attempt to convert to base unit
		newQty, newUnit, err := resolveConversionChain(ingID, nextQty, c.tUnit, targetUnit, db, visited)
		if err == nil {
			return newQty, newUnit, nil
		}
	}

	return 0, "", fmt.Errorf("cannot convert '%s' → '%s' for ingredient", fromUnit, targetUnit)
}
