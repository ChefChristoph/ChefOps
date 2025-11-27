package main

import (
	"fmt"
	"strings"

	"github.com/ChefChristoph/chefops/internal/tui"
)

func detailView(m Model) string {
	if m.activeRecipe == nil {
		return "Error: no recipe loaded"
	}

	r := m.activeRecipe

	// Load notes for this recipe
	notes, err := tui.LoadRecipeNotes(m.db, r.ID)
	if err != nil {
		notes = "Error loading notes"
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Recipe: %s\n", r.Name))
	builder.WriteString("──────────────────────────────────────\n")

	builder.WriteString(fmt.Sprintf("Yield: %.2f %s\n", r.YieldQty, r.YieldUnit))
	builder.WriteString(fmt.Sprintf("Total Cost: %.2f\n\n", r.TotalCost))

	builder.WriteString(fmt.Sprintf("%-12s %-24s %-10s %-6s %-10s\n", "Type", "Name", "Qty", "Unit", "Cost"))
	builder.WriteString("---------------------------------------------------------------------\n")

	for _, l := range r.Lines {
		builder.WriteString(fmt.Sprintf("%-12s %-24s %-10.3f %-6s %-10.2f\n",
			l.Type, l.Name, l.Qty, l.Unit, l.Cost))
	}

	// Add notes section
	builder.WriteString("\n────────── Notes ──────────\n")
	if notes == "" {
		builder.WriteString("(No notes available)")
	} else {
		// Truncate notes if too long for display
		maxLines := 10
		lines := strings.Split(notes, "\n")
		if len(lines) > maxLines {
			lines = lines[:maxLines]
			builder.WriteString(strings.Join(lines, "\n"))
			builder.WriteString(fmt.Sprintf("\n... (%d more lines)", len(lines)-maxLines))
		} else {
			builder.WriteString(notes)
		}
	}

	// Add export option
	builder.WriteString("\n\n")
	if m.detailCursor == 0 {
		builder.WriteString(activeItemStyle.Render("↳ Export this recipe to CSV"))
	} else {
		builder.WriteString("↳ Export this recipe to CSV")
	}

	builder.WriteString("\n")
	if m.detailCursor == 1 {
		builder.WriteString(activeItemStyle.Render("← Back"))
	} else {
		builder.WriteString("← Back")
	}

	builder.WriteString(" • ")
	if m.detailCursor == 2 {
		builder.WriteString(activeItemStyle.Render("q Quit"))
	} else {
		builder.WriteString("q Quit")
	}

	return detailStyle.Render(builder.String())
}

func exportConfirmView(m Model) string {
	var builder strings.Builder

	if m.exportError != "" {
		builder.WriteString("Error writing file:\n")
		builder.WriteString(m.exportError)
		builder.WriteString("\n\nPress Enter to return.")
	} else {
		builder.WriteString("✔ Exported successfully\n")
		builder.WriteString(m.exportPath)
		builder.WriteString("\n\nPress Enter to continue...")
	}

	// Center the content
	lines := strings.Split(builder.String(), "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	var centered strings.Builder
	for _, line := range lines {
		if line == "" {
			centered.WriteString("\n")
		} else {
			padding := (maxWidth - len(line)) / 2
			centered.WriteString(strings.Repeat(" ", padding) + line + "\n")
		}
	}

	return detailStyle.Render(centered.String())
}
