package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ChefChristoph/chefops/internal"
	"github.com/ChefChristoph/chefops/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	paneBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	listStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212"))

	activeItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("51")).
			Bold(true)

	detailStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
)

type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenRecipes
	ScreenMetadataImport
	ScreenNotesImport
)

type screen int

const (
	screenList screen = iota
	screenDetail
	screenExportConfirm
)

type Model struct {
	db *sql.DB

	currentScreen screen
	screen        Screen
	menuIndex     int

	// list screen
	recipes []tui.RecipeSummary
	cursor  int

	// detail screen
	activeRecipe *tui.RecipeDetail
	detailCursor int
	exportPath   string
	exportError  string

	// metadata import screen
	selectedRecipeID int
	selectedFile     string
	metaFiles        []string
	metaCursor       int

	// notes import screen
	noteFiles        []string
	noteCursor       int
	selectedNoteFile string
}

func NewModel(db *sql.DB) (*Model, error) {
	list, err := tui.LoadRecipes(db)
	if err != nil {
		return nil, err
	}

	return &Model{
		db:            db,
		recipes:       list,
		currentScreen: screenList,
		screen:        ScreenDashboard,
		menuIndex:     0,
	}, nil
}

func (m Model) Init() tea.Cmd { return nil }

func (m *Model) loadMetaFiles() {
	// Load files from ./recipe_meta/ directory
	m.metaFiles = []string{}

	// Read actual files from directory
	files, err := filepath.Glob("recipe_meta/*")
	if err == nil {
		for _, file := range files {
			// Only include files with supported extensions
			ext := strings.ToLower(filepath.Ext(file))
			if ext == ".md" || ext == ".json" || ext == ".txt" {
				m.metaFiles = append(m.metaFiles, file)
			}
		}
	}
}

func (m *Model) loadNoteFiles() {
	// Load files from both ./notes/ and ./recipe_notes/ directories
	m.noteFiles = []string{}

	// Check both directories
	directories := []string{"notes/*", "recipe_notes/*"}

	for _, pattern := range directories {
		files, err := filepath.Glob(pattern)
		if err == nil {
			for _, file := range files {
				// Only include markdown and text files
				ext := strings.ToLower(filepath.Ext(file))
				if ext == ".md" || ext == ".txt" {
					m.noteFiles = append(m.noteFiles, file)
				}
			}
		}
	}
}

func (m *Model) importMetadata() {
	if m.selectedRecipeID == 0 || m.selectedFile == "" {
		return
	}

	// Load existing metadata
	existingMeta, err := internal.LoadRecipeMetadata(m.selectedRecipeID)
	if err != nil {
		// Handle error - maybe set a message in the model
		return
	}

	// Load metadata from file
	newMeta, err := internal.LoadMetadataFromFile(m.selectedFile)
	if err != nil {
		// Handle error
		return
	}

	// Update timestamp
	internal.UpdateTimestamp(newMeta)

	// Merge metadata
	mergedMeta := internal.MergeMetadata(existingMeta, newMeta)

	// Save to database
	err = internal.SaveRecipeMetadata(m.selectedRecipeID, mergedMeta)
	if err != nil {
		// Handle error
		return
	}

	// Reset selection
	m.selectedRecipeID = 0
	m.selectedFile = ""
	m.screen = ScreenDashboard
}

func (m *Model) importNotes() {
	if m.selectedNoteFile == "" {
		return
	}

	// Extract recipe name from filename
	filename := filepath.Base(m.selectedNoteFile)
	recipeName := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Handle different filename patterns
	// Remove prefixes like "DISH ", "BULK ", etc.
	cleanName := recipeName
	if strings.HasPrefix(cleanName, "DISH ") {
		cleanName = strings.TrimPrefix(cleanName, "DISH ")
	} else if strings.HasPrefix(cleanName, "BULK ") {
		cleanName = strings.TrimPrefix(cleanName, "BULK ")
	}

	// Try to find recipe by name (try both original and cleaned)
	recipeID, err := internal.GetRecipeIDByName(cleanName)
	if err != nil {
		// Try with original name if cleaned didn't work
		recipeID, err = internal.GetRecipeIDByName(recipeName)
		if err != nil {
			// Recipe not found, could add manual selection here
			return
		}
	}

	// Load and process notes from file
	notes, err := internal.LoadNotesFromFile(m.selectedNoteFile)
	if err != nil {
		// Handle error
		return
	}

	// Update database
	err = internal.UpdateRecipeNotes(m.db, recipeID, notes)
	if err != nil {
		// Handle error
		return
	}

	// Reset selection
	m.selectedNoteFile = ""
	m.screen = ScreenDashboard
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// universal quit
		case "q", "ctrl+c":
			return m, tea.Quit

		// DASHBOARD SCREEN INPUT
		case "up", "k":
			if m.screen == ScreenDashboard && m.menuIndex > 0 {
				m.menuIndex--
			} else if m.currentScreen == screenList && m.cursor > 0 {
				m.cursor--
			} else if m.currentScreen == screenDetail && m.detailCursor > 0 {
				m.detailCursor--
			} else if m.screen == ScreenMetadataImport {
				if m.selectedRecipeID == 0 && m.cursor > 0 {
					m.cursor--
				} else if m.selectedRecipeID != 0 && m.metaCursor > 0 {
					m.metaCursor--
				}
			} else if m.screen == ScreenNotesImport && m.noteCursor > 0 {
				m.noteCursor--
			}
		case "down", "j":
			if m.screen == ScreenDashboard && m.menuIndex < len(dashboardItems)-1 {
				m.menuIndex++
			} else if m.currentScreen == screenList && m.cursor < len(m.recipes)-1 {
				m.cursor++
			} else if m.currentScreen == screenDetail && m.detailCursor < 2 {
				m.detailCursor++
			} else if m.screen == ScreenMetadataImport {
				if m.selectedRecipeID == 0 && m.cursor < len(m.recipes)-1 {
					m.cursor++
				} else if m.selectedRecipeID != 0 && m.metaCursor < len(m.metaFiles)-1 {
					m.metaCursor++
				}
			} else if m.screen == ScreenNotesImport && m.noteCursor < len(m.noteFiles)-1 {
				m.noteCursor++
			}
		case "enter":
			if m.screen == ScreenDashboard {
				switch m.menuIndex {
				case 0:
					m.screen = ScreenRecipes
				case 1:
					// Initialize metadata import screen
					m.screen = ScreenMetadataImport
					m.loadMetaFiles()
					m.metaCursor = 0
				case 2:
					// Initialize notes import screen
					m.screen = ScreenNotesImport
					m.loadNoteFiles()
					m.noteCursor = 0
				case 6:
					return m, tea.Quit
				}
			} else if m.currentScreen == screenList {
				// enter opens detail view
				id := m.recipes[m.cursor].ID
				detail, _ := tui.LoadRecipeDetail(m.db, id)
				m.activeRecipe = detail
				m.currentScreen = screenDetail
				m.detailCursor = 0
			} else if m.currentScreen == screenDetail {
				// Handle detail view menu selection
				switch m.detailCursor {
				case 0: // Export
					if m.activeRecipe != nil {
						// Generate export path
						slug := tui.Slugify(m.activeRecipe.Name)
						m.exportPath = fmt.Sprintf("exports/recipes/%s.csv", slug)
						m.exportError = ""

						// Export the recipe
						err := tui.ExportRecipeToCSV(m.activeRecipe, m.exportPath)
						if err != nil {
							m.exportError = err.Error()
						}

						// Show confirmation
						m.currentScreen = screenExportConfirm
					}
				case 1: // Back
					m.currentScreen = screenList
				case 2: // Quit
					return m, tea.Quit
				}
			} else if m.currentScreen == screenExportConfirm {
				// Return from export confirmation
				m.currentScreen = screenDetail
			} else if m.screen == ScreenMetadataImport {
				if m.selectedRecipeID == 0 {
					// Select recipe
					m.selectedRecipeID = m.recipes[m.cursor].ID
				} else if m.selectedFile == "" {
					// Select file and import
					if len(m.metaFiles) > 0 {
						m.selectedFile = m.metaFiles[m.metaCursor]
						m.importMetadata()
					}
				}
			} else if m.screen == ScreenNotesImport {
				if len(m.noteFiles) > 0 {
					m.selectedNoteFile = m.noteFiles[m.noteCursor]
					m.importNotes()
				}
			}
		case "b":
			if m.screen == ScreenRecipes || m.screen == ScreenMetadataImport || m.screen == ScreenNotesImport {
				m.screen = ScreenDashboard
				m.selectedRecipeID = 0
				m.selectedFile = ""
				m.selectedNoteFile = ""
			}
		case "esc", "backspace", "left", "h":
			if m.currentScreen == screenDetail || m.currentScreen == screenExportConfirm {
				m.currentScreen = screenList
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.screen {
	case ScreenDashboard:
		return renderDashboard(m)
	case ScreenRecipes:
		if m.currentScreen == screenExportConfirm {
			// Show export confirmation as full screen
			return exportConfirmView(m)
		}
		left := listView(m)
		right := detailView(m)

		leftPane := paneBorder.Width(40).Render(left)
		rightPane := paneBorder.Width(70).Render(right)

		ui := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

		footer := footerStyle.Render("\n↑/↓ navigate • Enter open • b back • q quit")

		return ui + footer
	case ScreenMetadataImport:
		return renderMetadataImport(m)
	case ScreenNotesImport:
		return renderNotesImport(m)
	default:
		return "unknown screen"
	}
}
