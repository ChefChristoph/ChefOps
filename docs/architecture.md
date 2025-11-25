# ChefOps Architecture

## Overview
ChefOps is a lightweight command-line tool for chefs and kitchen operators.  
It is built around three core components:

1. **CLI (Go)**  
   Handles commands like ingredient management, recipe costing, scaling, and forecasting.

2. **SQLite Database**  
   Stores ingredients, recipes, recipe-item relations, and subrecipes.

3. **SQL Views**  
   Provide high-level abstractions for:
   - Cost expansion
   - Subrecipe nesting
   - Ingredient rollups
   - Marketlist generation

---

## Component Diagram
+————————+
|        chefops         |
|     (Go CLI app)       |
+———–+————+
|
v
+————————+
|     SQLite Database    |
+———–+————+
|
v
+————————+
|      SQL Views         |
| recipe_raw_lines       |
| recipe_items_expanded  |
| recipe_items_expanded_detail |
| recipe_totals          |
| market_list            |
+————————+

---

## Data Flow Example: Forecast

1. User calls:  
   `chefops forecast "DISH Pole Position Burger" --portions 150`

2. CLI resolves recipe ID → loads raw lines

3. SQL expansion generates:
   - Direct ingredient quantities  
   - Subrecipe-level requirements  
   - Total quantities

4. Output:  
   Scaled, sorted, human-readable production plan.

---

## Design Principles

- **Chef-first:** The CLI never forces complex flags or computer terminology.
- **Predictable:** All quantities normalized to SI units (kg, liter, piece).
- **Transparent:** All costing visible down to line level.
- **Composable:** Recipes can reference other recipes.
- **Reproducible:** Master import script can rebuild the entire DB anytime.


