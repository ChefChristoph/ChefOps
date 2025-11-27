package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ChefChristoph/chefops/internal"
)

func main() {
	db, err := internal.OpenDB()
	if err != nil {
		fmt.Println("Failed to open DB:", err)
		os.Exit(1)
	}
	defer db.Close()

	m, err := NewModel(db)
	if err != nil {
		fmt.Println("Failed to load recipes:", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}
