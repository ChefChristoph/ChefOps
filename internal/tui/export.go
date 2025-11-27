package tui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Slugify converts a recipe name to a filename-safe format
func Slugify(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces with underscores
	slug = strings.ReplaceAll(slug, " ", "_")

	// Remove special characters except underscores and hyphens
	reg := regexp.MustCompile(`[^a-z0-9_]`)
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive underscores
	reg = regexp.MustCompile(`_+`)
	slug = reg.ReplaceAllString(slug, "_")

	// Remove leading/trailing underscores
	slug = strings.Trim(slug, "_")

	return slug
}

// ExportRecipeToCSV exports a recipe to CSV format
func ExportRecipeToCSV(recipe *RecipeDetail, path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Create file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write header comments
	fmt.Fprintf(writer, "Recipe: %s\n", recipe.Name)
	fmt.Fprintf(writer, "Yield: %.2f %s\n", recipe.YieldQty, recipe.YieldUnit)
	fmt.Fprintf(writer, "Total Cost: %.2f\n", recipe.TotalCost)
	fmt.Fprintf(writer, "\n")

	// Write CSV header
	fmt.Fprintf(writer, "Type,Name,Qty,Unit,Cost\n")

	// Write recipe lines
	for _, line := range recipe.Lines {
		fmt.Fprintf(writer, "%s,%s,%.3f,%s,%.2f\n",
			line.Type, line.Name, line.Qty, line.Unit, line.Cost)
	}

	return nil
}
