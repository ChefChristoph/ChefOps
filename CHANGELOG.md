# Changelog
All notable changes to **ChefOps** will be documented in this file.

This project follows the structure recommended by [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and uses semantic versioning (MAJOR.MINOR.PATCH).

---

## [0.6.0] – 2025-11-26
### Added
- **Recipe Export to CSV in TUI** (opencode)
  - New "Export this recipe to CSV" option in recipe detail view
  - Interactive navigation with arrow keys between Export, Back, and Quit options
  - CSV export to `exports/recipes/<slugified-name>.csv` with complete recipe data
  - Export confirmation dialog showing success/error messages with file paths

- **CSV Export Engine** (`internal/tui/export.go`)
  - `Slugify()` function to convert recipe names to filename-safe format
    - Example: "DISH Hot Lap Honey & Pepperoni Pizza" → `dish_hot_lap_honey_pepperoni_pizza.csv`
  - `ExportRecipeToCSV()` function generating properly formatted CSV files
  - Automatic directory creation for `exports/recipes/` folder
  - CSV format includes recipe metadata (name, yield, total cost) and ingredient breakdown

- **Enhanced TUI Detail View**
  - Made detail view interactive with selectable menu items
  - Added cursor navigation for export options
  - New `screenExportConfirm` state for export confirmation dialogs
  - Improved error handling with user-friendly error messages

### Changed
- Updated TUI navigation to support interactive detail view menu
- Enhanced Model struct with export-related fields (`detailCursor`, `exportPath`, `exportError`)
- Modified detail view rendering to show active selection state
- Improved TUI state management for export workflow

### Fixed
- Missing import for `fmt` package in model.go
- Enhanced navigation logic to handle detail view cursor movement
- Proper state transitions between detail view and export confirmation

### Tested
- ✅ CSV export functionality with real recipe data
- ✅ Slugify function with various recipe name formats (DISH, BULK, special characters)
- ✅ Interactive TUI navigation and selection
- ✅ Export confirmation dialogs (success and error cases)
- ✅ File creation and directory structure (`exports/recipes/`)
- ✅ CSV output format matching specification (Type,Name,Qty,Unit,Cost)

---

## [0.5.0] – 2025-11-26
### Added
- **Enhanced Notes Import System** (opencode)
  - Multi-format support: `.md`, `.txt`, and `.json` files for notes import
  - Dual directory support: searches both `./notes/` and `./recipe_notes/` directories
  - Smart recipe name matching with prefix handling (`DISH `, `BULK `)
  - JSON notes extraction from `notes` field in structured JSON files
  - Enhanced file picker interface with improved file discovery

- **Improved TUI Notes Workflow**
  - Updated `loadNoteFiles()` function to scan multiple directories
  - Better filename parsing and recipe name resolution
  - Enhanced error handling for missing recipes and files
  - Improved user feedback during import process

- **Database Integration**
  - Added `LoadRecipeNotes()` function to `internal/tui/db.go`
  - Fixed import paths in detail view for proper notes loading
  - Enhanced notes storage and retrieval with proper error handling

### Changed
- Updated notes import workflow to handle multiple file formats and directories
- Improved recipe name extraction from filenames with common prefixes
- Enhanced file picker to show files from both `notes/` and `recipe_notes/` directories
- Updated detail view imports to use correct package references

### Fixed
- Missing `LoadRecipeNotes` function in `internal/tui/db.go`
- Incorrect import path in `cmd/tui/detail_view.go`
- Limited directory support in notes file picker (now supports both directories)
- Recipe name matching issues with prefixed filenames

### Tested
- ✅ Multi-format file parsing (Markdown, Text, JSON)
- ✅ Dual directory file discovery
- ✅ Recipe name matching with prefix handling
- ✅ Notes extraction filtering tables and ingredients
- ✅ Database storage and retrieval operations
- ✅ TUI integration and user workflow

---

## [0.4.0] – 2025-11-25
### Added
- **Markdown → Recipe Notes Importer** (opencode)
  - New `recipe_notes/` directory for markdown files containing recipe instructions and notes
  - Smart markdown parser that extracts freeform text while filtering out structured data
  - Safe parsing that preserves operational knowledge without interfering with recipe system

- **CLI notes commands**
  - `chefops recipe note import --recipe "Recipe Name" --file path/to/file.md` - Import notes from markdown
  - `chefops recipe note show "Recipe Name"` - Display stored notes for a recipe
  - Support for recipe names with spaces and complex file paths
  - Colored success/error feedback using lipgloss styling

- **TUI notes integration**
  - New dashboard item: "📝 Import Notes From File"
  - File picker interface for `./recipe_notes/` directory
  - Automatic recipe name mapping from filename
  - Notes viewer in recipe detail view with scroll indication
  - Success feedback and error handling in TUI interface

- **Smart markdown parser** (`internal/notes.go`)
  - `ExtractNotesFromMarkdown()` function that safely filters content:
    - ✅ **Extracts**: Instructions, notes, tips, headings, paragraphs, YAML frontmatter
    - ❌ **Filters**: Markdown tables (`| col | col |`), table separators (`---|---`)
    - ❌ **Filters**: Ingredient-like lines with quantities (`- 500g flour`, `* 2 kg sugar`)
  - Regex-based ingredient detection for common units (kg, g, l, ml, piece, etc.)
  - Preserves markdown formatting and structure for freeform text

- **Database notes functions**
  - `UpdateRecipeNotes(db, recipeID, notes)` - Save notes to recipes.notes column
  - `LoadRecipeNotes(db, recipeID)` - Load notes from recipes.notes column
  - `LoadNotesFromFile(filepath)` - Read and process markdown files for import
  - NULL-safe handling for existing recipes without notes

- **Enhanced TUI detail view**
  - Added "────────── Notes ──────────" section to recipe detail view
  - Displays up to 10 lines of notes with overflow indication
  - Shows "(No notes available)" when no notes present
  - Integrated with existing recipe information display

### Changed
- Updated `RecipeDetail` struct to include `ID` field for notes loading
- Enhanced CLI argument parsing for commands with `--recipe` and `--file` flags
- Modified TUI navigation to handle new notes import screen
- Updated dashboard menu to include notes import option

### Fixed
- Import variable redeclaration conflicts in CLI commands
- Missing function definitions for TUI notes import workflow
- Recipe ID resolution for notes loading in detail view
- File path handling for names with spaces

### Tested
- ✅ Real-world markdown files (Bulk_Pizza_Sauce.md, Bulk_Batter.md, etc.)
- ✅ Complex markdown with YAML frontmatter, tables, and ingredient lists
- ✅ Parser correctly filters tables and ingredient-like lines
- ✅ Preserves instructions, tips, storage notes, and quality indicators
- ✅ Maintains existing functionality (recipe list, show, cost, scaling)

---

## [0.3.0] – 2025-11-25
### Added
- **Recipe metadata system** (opencode)
  - New `metadata TEXT` column in recipes table for storing JSON metadata
  - Complete metadata struct with fields: description, instructions, notes, mise_en_place, allergens, equipment, tags, created_by, last_updated
  - Smart metadata merging that preserves existing fields when updating
  - Automatic timestamp updates on metadata changes

- **CLI metadata commands**
  - `chefops recipe set-meta "Recipe Name" filepath` - Import metadata from files
  - `chefops recipe export-meta "Recipe Name" --format=json|md` - Export metadata
  - Support for multiple file formats:
    - `.md` - Parse markdown sections into metadata fields
    - `.json` - Direct JSON import/export
    - `.txt` - Add content to notes section
  - Proper handling of recipe names with spaces

- **TUI metadata import**
  - New menu item: "📥 Import / Update Recipe Metadata"
  - Two-step workflow: select recipe → select file → import
  - Reads actual files from `./recipe_meta/` directory
  - Success feedback and error handling in TUI interface

- **Metadata parsing engine**
  - Markdown parser supporting sections: Description, Instructions, Notes, Mise En Place, Allergens, Equipment, Tags, Created By, Last Updated
  - List item parsing for bullet points (- item, * item)
  - JSON serialization/deserialization with proper error handling
  - Text file handling for simple note additions

- **Database helper functions**
  - `LoadRecipeMetadata(recipeID)` - Load and parse metadata from database
  - `SaveRecipeMetadata(recipeID, metadata)` - Save metadata as JSON
  - `GetRecipeIDByName(name)` - Helper to resolve recipe names to IDs
  - NULL-safe metadata handling for existing recipes

### Changed
- Updated schema.sql to include metadata column for new installations
- Enhanced CLI argument parsing for complex commands with flags and spaced arguments
- Improved error handling with colored success/error messages using lipgloss

### Fixed
- Database NULL handling for metadata column when loading existing recipes
- Argument parsing for export commands with --format flag
- Recipe name resolution for names containing spaces

---

## [0.2.2] – 2025-11-25
### Changed
- **TUI UI improvements** (opencode)
  - Enhanced visual layout with proper lipgloss styling
  - Implemented side-by-side pane layout using `lipgloss.JoinHorizontal`
  - Added global styles for consistent theming across components
  - Recipe list now uses color-coded active items with bold highlighting
  - Detail view wrapped in consistent styling
  - Footer styling applied for navigation hints
  - Removed redundant navigation text from individual panes

### Fixed
- **TUI component styling**
  - Centralized style definitions in `view.go` to prevent duplication
  - List items now properly styled with `listStyle` and `activeItemStyle`
  - Detail view content wrapped with `detailStyle`
  - Error messages in detail view properly styled

---

## [0.2.1] – 2025-11-25
### Fixed
- **Export functionality** (opencode Zen Big Pickle)
  - Fixed column name mismatches in marketlist export (`total_qty_needed` → `total_qty`, `estimated_cost` → `total_cost`)
  - Created missing `recipe_items_expanded_detail_export` view for recipe export functionality
  - Updated SQL queries in export commands to use correct view structure
  - Applied database views to ensure export functionality works correctly

- **Export commands now fully operational**
  - `chefops export marketlist` - Working with both markdown and JSON output
  - `chefops export recipe "NAME"` - Individual recipes with full ingredient breakdown and costing
  - `chefops export full-report` - Complete recipe catalog with yields
  - File output with `-o` flag working correctly
  - JSON format with `--json` flag producing structured data

---

## [0.2.0] – 2025-11-25
### Added
- **Full recipe costing engine**
  - Ingredient-level cost calculation  
  - Subrecipe cost expansion based on yield  
  - `chefops recipe cost NAME` command  

- **Recipe scaling**
  - New command:  
    ```sh
    chefops recipe scale "NAME" --qty X --unit UNIT
    ```  
  - Supports scaling to any yield in any unit  

- **Forecasting engine**
  - New command:  
    ```sh
    chefops forecast "DISH NAME" --portions X
    ```  
  - Outputs:  
    - Required subrecipes (scaled quantities)  
    - Required ingredients  
    - Units preserved automatically  
  - Useful for production planning, events, and procurement  

- **New database views**
  - `recipe_raw_lines`
  - `recipe_items_expanded_detail`
  - Improved `recipe_totals`
  - Improved `market_list`

- **Marketlist generator**
  - `chefops marketlist`  
  - Summarizes total ingredients required across all recipes

### Changed
- Reworked schema to add:
  - `unit` column for `recipe_subrecipes`
  - Better yield handling for recipes  
- Rebuilt the entire import pipeline (`MASTER_import_F1.sh`)
- Improved CLI argument parsing for commands with names containing spaces

### Fixed
- Subrecipe unit inconsistencies  
- Recursive cost calculation  
- Circular view definitions  
- “ingredient not found” fuzzy matcher improvements  
- Multiple schema inconsistencies  
- Corrected recipe show display  
- Corrected view reloading & resetting logic  
- Resolved case-sensitivity issues between local FS and GitHub  

---

## [0.1.0] – Initial Development Snapshot
### Added
- Basic project structure  
- Ingredient import + listing  
- Simple recipe creation & ingredient linking  
- Early experimental schema (pre-refactor)

---

## Upcoming
### Planned Features
- Rich recipe metadata editing in TUI (not just import)
- Metadata search and filtering capabilities
- Recipe metadata templates for common dish types
- Enhanced notes viewer with scrolling and search
- Batch notes import from entire directory
- Portion-based cost menu engineering
- Export to Google Sheets / Excel
- Calendar-based production planning
- Inventory integration  
- Supplier pricing histories  
- Automated batch scaling per forecast