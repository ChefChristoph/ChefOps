# ChefOps – TUI Dashboard Roadmap (Bubble Tea)

A safe, gradual plan to build the ChefOps Terminal UI without breaking the existing CLI.

---

## 0. Safety & Setup
- [ ] Create a new branch: `feature/tui-dashboard`
- [ ] Ensure current CLI builds cleanly (`go build -o chefops ./cmd/chefops`)
- [ ] Tag main as `v0.2.0` (already done)
- [ ] Create TUI as a **separate entry point**: `cmd/chefops-tui/`
- [ ] Do NOT modify any existing CLI files

---

## 1. Basic TUI Skeleton (Bubble Tea)
Goal: A basic UI window with a placeholder.

- [ ] Create `cmd/chefops-tui/main.go`
- [ ] Add Bubble Tea, Lipgloss, Bubbles dependencies
- [ ] Render a simple “ChefOps TUI Loaded” screen
- [ ] Add ESC/Q to quit
- [ ] Add Makefile entry: `make tui`

**Result:**  
TUI runs with:  
go run ./cmd/chefops-tui

---

## 2. Load Recipes From SQLite
Goal: Show recipe list on the left.

- [ ] Create DB loader in TUI folder
- [ ] Reuse your `internal.OpenDB()`
- [ ] Query recipes: `id, name, yield_qty, yield_unit`
- [ ] Create a Bubbletea list component
- [ ] Add arrow-key navigation

**Result:**  
Left side shows all recipes, scrollable.

---

## 3. Recipe Detail Panel
Goal: Show ingredient + subrecipe breakdown for selected recipe.

- [ ] Query `recipe_raw_lines`
- [ ] Display table with:
  - type  
  - qty  
  - unit  
  - name  
  - line cost  
- [ ] Add Lipgloss styling
- [ ] Auto-refresh when selection changes

---

## 4. Tabs System
Goal: Multi-section dashboard

Tabs:
1. Recipes  
2. Ingredients  
3. Marketlist  
4. Scaling  
5. Forecast  
6. Export  

- [ ] Add “Tabs” bubble component
- [ ] Add Tab switching with ←→ or 1–6
- [ ] Render relevant view on right side

---

## 5. Scaling Popup
Goal: Perform scaling inside the TUI

- [ ] Press "s" to open modal  
- [ ] Enter qty + unit  
- [ ] Show scaled table  
- [ ] Export to CSV

Uses your existing scaling logic.

---

## 6. Forecast Popup
Goal: Forecast items interactively

- [ ] Press "f" to open modal  
- [ ] Input portions  
- [ ] Show aggregated ingredient list  
- [ ] Save to CSV

Uses your forecast engine.

---

## 7. Marketlist View
- [ ] Display current marketlist view from database
- [ ] Add filter/search
- [ ] Export to CSV

---

## 8. Export Menu
- [ ] Press "e" → choose CSV, JSON, Google Sheet
- [ ] Use existing exporters
- [ ] Google Sheets: optional later

---

## 9. Polishing & Theming
- [ ] Lipgloss theme
- [ ] Headers, borders
- [ ] Status bar
- [ ] Error toasts
- [ ] Loading spinners
- [ ] Keybindings help panel

---

## 10. Merge & Release
- [ ] Merge branch once tested  
- [ ] Tag version `v0.3.0-TUI`  
- [ ] Update README  
- [ ] Add a screenshot to repo  
- [ ] Share release on GitHub  

---

# Notes
- Existing CLI **remains untouched and fully functional**.
- TUI uses the same database and functions but is a *separate app*.
- You can keep building recipes, scaling, forecasts, OSS, F1 planning normally.

---

# Optional Future Ideas
- [ ] Apple Notes export
- [ ] n8n integration
- [ ] Recipe editing inside TUI
- [ ] Ingredient master editor
- [ ] Live price updates via API
- [ ] AI helper panel (local Ollama)


# Example screen:
┌──────────────────────────────────────────────────────────────────────────────┐
│                                 CHEF OPS TUI                                 │
├──────────────────────────────────────────────────────────────────────────────┤
│ SEARCH:  pole burger                                                       ⌕ │
├───────────────┬──────────────────────────────────────────────────────────────┤
│  Recipes      │  DISH Pole Position Burger                                    │
│  Ingredients  │  ───────────────────────────────────────────────────────────   │
│  Costs        │   Beef Patty...................... 0.180 kg      5.76 AED     │
│  Scaling      │   Cheddar Cheese.................. 0.020 kg      0.52 AED     │
│  Forecast     │   Bun (Burger).................... 1 piece       1.50 AED     │
│  Marketlist   │   Pickles......................... 0.020 kg      0.20 AED     │
│               │   Tomato.......................... 0.030 kg      0.15 AED     │
│               │   SUB Mac & Cheese Base........... 0.060 kg      0.66 AED     │
│               │   BULK Burger Sauce............... 0.030 kg      0.30 AED     │
│               │   TOTAL COST..................................... 10.26 AED    │
├───────────────┴──────────────────────────────────────────────────────────────┤
│  ↑↓ Move   Enter Inspect   S Scale   F Forecast   M Marketlist   Q Quit        │
└──────────────────────────────────────────────────────────────────────────────┘

