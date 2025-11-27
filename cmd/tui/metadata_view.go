package main

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

var (
	metaTitleStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

	metaBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	metaItemStyle = lipgloss.NewStyle()
	metaSelected  = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
)

func renderMetadataImport(m Model) string {
	var lines []string

	lines = append(lines, metaTitleStyle.Render("📥 Import / Update Recipe Metadata"))
	lines = append(lines, "Use ↑/↓ and Enter to select, b to go back, q to quit.")
	lines = append(lines, "")

	if m.selectedRecipeID == 0 {
		// Recipe selection
		lines = append(lines, "Step 1: Select a recipe")
		lines = append(lines, "")

		for i, recipe := range m.recipes {
			prefix := "  "
			style := metaItemStyle
			if i == m.cursor {
				prefix = "➤ "
				style = metaSelected
			}
			lines = append(lines, style.Render(fmt.Sprintf("%s%s", prefix, recipe.Name)))
		}
	} else if m.selectedFile == "" {
		// File selection
		selectedRecipe := "Unknown"
		for _, recipe := range m.recipes {
			if recipe.ID == m.selectedRecipeID {
				selectedRecipe = recipe.Name
				break
			}
		}

		lines = append(lines, fmt.Sprintf("Step 2: Select metadata file for recipe: %s", selectedRecipe))
		lines = append(lines, "")

		if len(m.metaFiles) == 0 {
			lines = append(lines, "No metadata files found in ./recipe_meta/ directory")
		} else {
			for i, file := range m.metaFiles {
				prefix := "  "
				style := metaItemStyle
				if i == m.metaCursor {
					prefix = "➤ "
					style = metaSelected
				}
				filename := filepath.Base(file)
				lines = append(lines, style.Render(fmt.Sprintf("%s%s", prefix, filename)))
			}
		}
	} else {
		// Importing
		lines = append(lines, fmt.Sprintf("Importing metadata from: %s", filepath.Base(m.selectedFile)))
		lines = append(lines, "Processing...")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return metaBoxStyle.Render(content)
}
