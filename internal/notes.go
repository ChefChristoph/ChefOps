package internal

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ExtractNotesFromMarkdown extracts freeform text from markdown, filtering out tables and ingredient lists
func ExtractNotesFromMarkdown(md string) string {
	if md == "" {
		return ""
	}

	var result []string
	scanner := bufio.NewScanner(strings.NewReader(md))

	// Regex to detect ingredient-like lines with quantities and units
	ingredientRegex := regexp.MustCompile(`^\s*[-*]?\s*\d+(\.\d+)?\s*(kg|g|l|ml|piece|pc|lb|oz|cup|tbsp|tsp|tablespoon|teaspoon)`)

	// Track if we're in a table
	inTable := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Skip table markers
		if strings.HasPrefix(trimmed, "|") || strings.HasPrefix(trimmed, "|-") || strings.Contains(trimmed, "---|---") {
			inTable = true
			continue
		}

		// If we were in a table and this line doesn't start with |, we're out of the table
		if inTable && !strings.HasPrefix(trimmed, "|") {
			inTable = false
		}

		// Skip lines while in table
		if inTable {
			continue
		}

		// Skip ingredient-like lines
		if ingredientRegex.MatchString(trimmed) {
			continue
		}

		// Skip empty lines at the start, but preserve them for structure
		if len(result) == 0 && trimmed == "" {
			continue
		}

		// Keep all other lines
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// UpdateRecipeNotes updates the notes field for a recipe
func UpdateRecipeNotes(db *sql.DB, recipeID int, notes string) error {
	query := "UPDATE recipes SET notes = ? WHERE id = ?"
	_, err := db.Exec(query, notes, recipeID)
	if err != nil {
		return fmt.Errorf("failed to update recipe notes: %w", err)
	}
	return nil
}

// LoadRecipeNotes loads the notes field for a recipe
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

// LoadNotesFromFile reads and processes a file for notes
func LoadNotesFromFile(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	ext := strings.ToLower(filepath[strings.LastIndex(filepath, ".")+1:])

	switch ext {
	case "json":
		// Parse JSON and extract notes field
		var jsonData map[string]interface{}
		err := json.Unmarshal(data, &jsonData)
		if err != nil {
			return "", fmt.Errorf("failed to parse JSON: %w", err)
		}

		if notes, ok := jsonData["notes"].(string); ok {
			return ExtractNotesFromMarkdown(notes), nil
		}
		return "", fmt.Errorf("no notes field found in JSON")

	case "md", "txt":
		content := string(data)
		cleaned := ExtractNotesFromMarkdown(content)
		return cleaned, nil

	default:
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}
}
