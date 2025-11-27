package main

import "strings"

func listView(m Model) string {
	var b strings.Builder

	b.WriteString("ChefOps Recipe Browser\n")
	b.WriteString("──────────────────────────\n\n")

	for i, recipe := range m.recipes {
		if i == m.cursor {
			b.WriteString(activeItemStyle.Render("> " + recipe.Name))
		} else {
			b.WriteString(listStyle.Render("  " + recipe.Name))
		}
		b.WriteString("\n")
	}

	return b.String()
}
