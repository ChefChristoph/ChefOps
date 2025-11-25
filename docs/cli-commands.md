# ChefOps CLI Commands

## Ingredient Commands
### Add ingredient
chefops ingredient add --name "Tomato" --unit kg --cost 5

### List ingredients
chefops ingredient list

## Recipe Commands

### Create recipe

chefops recipe new --name "BULK Pasta Base" --yield 1 --unit kg
### Add ingredient
chefops recipe add-item --recipe "BULK Pasta Base" --ingredient "Butter" --qty 0.02
### Add subrecipe
chefops recipe add-subrecipe --recipe "Dish" --sub "BULK Base" --qty 0.1 --unit kg
### Show recipe
chefops recipe show "BULK Pizza Sauce"
### Cost recipe
chefops recipe cost "DISH Lobster Roll"
### Scale recipe
chefops recipe scale "BULK Batter" --qty 10 --unit kg

## Forecasting

Calculate ingredients for X portions:
chefops forecast "DISH Turbo Hammour Popcorn" --portions 150
Outputs:
	•	Required subrecipes (scaled)
	•	Required ingredients
	•	Marketlist-compatible totals

---

# 📄 **docs/import-pipeline.md**
```markdown
# Import Pipeline (MASTER_import_F1.sh)

The import pipeline performs:

1. **Deletes and recreates the database**
2. **Loads schema and views**
3. **Imports all ingredients with pricing**
4. **Imports all bulk recipes**
5. **Imports all dish recipes**
6. **Verifies costing**
7. Displays completion message

---

## Reset + Rebuild
```sh
rm -rf db
mkdir db
sqlite3 db/chefops.db < schema.sql
sqlite3 db/chefops.db < views.sql

Ingredient Import

Each ingredient is added with:
	•	name
	•	unit
	•	cost per unit

Example:

./chefops ingredient add --name "Lemon" --unit kg --cost 7

Bulk Recipe Import

Each bulk recipe:
	•	Declares base yield
	•	Adds ingredients or subrecipes
	•	Validates relations

⸻

Dish Recipe Import

Links:
	•	Ingredients
	•	Bulk preps
	•	Subrecipes
	•	Base yields

⸻

Structure

Your full F1 dataset is cleanly rebuildable from this file at any time.

---

# 📄 **docs/roadmap.md**
```markdown
# ChefOps Roadmap

## Near-Term
- Unit conversion engine (g ↔ kg, ml ↔ liter)
- Export to Excel / Google Sheets
- Round quantities by unit type
- Error logs + debug mode

## Mid-Term
- Inventory + supplier integration
- Costing history tracking
- Menu engineering reports
- Production batch planning

## Long-Term
- Full ChefOps TUI (terminal UI)
- Apple Shortcuts integration
- ChefOps Cloud Sync
- Team mode with shared DB
- iOS version (local SQLite)

---

Feel free to add features as your workflow expands.


