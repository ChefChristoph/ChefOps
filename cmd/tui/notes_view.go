package main

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

var (
	notesTitleStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

	notesBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	notesItemStyle = lipgloss.NewStyle()
	notesSelected  = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
)

func renderNotesImport(m Model) string {
	var lines []string

	lines = append(lines, notesTitleStyle.Render("📝 Import Notes From File"))
	lines = append(lines, "Use ↑/↓ and Enter to select, b to go back, q to quit.")
	lines = append(lines, "")
	lines = append(lines, "Available note files:")

	if len(m.noteFiles) == 0 {
		lines = append(lines, "No markdown files found in ./recipe_notes/ directory")
	} else {
		for i, file := range m.noteFiles {
			prefix := "  "
			style := notesItemStyle
			if i == m.noteCursor {
				prefix = "➤ "
				style = notesSelected
			}
			filename := filepath.Base(file)
			lines = append(lines, style.Render(fmt.Sprintf("%s%s", prefix, filename)))
		}
	}

	if m.selectedNoteFile != "" {
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("Importing notes from: %s", filepath.Base(m.selectedNoteFile)))
		lines = append(lines, "Processing...")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return notesBoxStyle.Render(content)
}
