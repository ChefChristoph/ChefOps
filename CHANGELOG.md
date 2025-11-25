# Changelog
All notable changes to **ChefOps** will be documented in this file.

This project follows the structure recommended by [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and uses semantic versioning (MAJOR.MINOR.PATCH).

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
- Portion-based cost menu engineering
- Export to Google Sheets / Excel
- Calendar-based production planning
- Inventory integration  
- Supplier pricing histories  
- Automated batch scaling per forecast