package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	dashboardTitleStyle = lipgloss.NewStyle().
				Bold(true).
				MarginBottom(1)

	dashboardBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(1, 2)

	dashboardItemStyle = lipgloss.NewStyle()
	dashboardSelected  = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
)

var dashboardItems = []string{
	"Browse recipes",
	"📥 Import / Update Recipe Metadata",
	"📝 Import Notes From File",
	"Forecast & scaling",
	"Market list",
	"Export to CSV",
	"Quit",
}

func renderDashboard(m Model) string {
	var lines []string

	lines = append(lines, dashboardTitleStyle.Render("ChefOps Dashboard"))
	lines = append(lines, "Use ↑/↓ and Enter, q to quit.")
	lines = append(lines, "")

	for i, item := range dashboardItems {
		prefix := "  "
		style := dashboardItemStyle
		if i == m.menuIndex {
			prefix = "➤ "
			style = dashboardSelected
		}
		lines = append(lines, style.Render(fmt.Sprintf("%s%s", prefix, item)))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return dashboardBoxStyle.Render(content)
}
