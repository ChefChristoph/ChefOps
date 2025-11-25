# SQL Views Documentation

ChefOps uses multiple SQL views to produce derived data.

---

## 1. `recipe_raw_lines`
Expands each recipe into:
- Direct ingredients
- Subrecipe references
- Line costs

Useful for:
- Recipe display
- Cost calculation

---

## 2. `recipe_items_expanded`
Aggregates all ingredient contributions across nested subrecipes.

Output:
| recipe_id | ingredient_id | total_qty |

---

## 3. `recipe_items_expanded_detail`
Detailed ingredient list with:
- Type (ingredient / subrecipe)
- Normalized unit
- Line quantities

Used for:
- `recipe show`
- Scaling
- Forecasting

---

## 4. `recipe_totals`
Produces:
- Total cost per recipe
- Cost per primary & secondary yield


Columns:

recipe_name
yield_qty
yield_unit
secondary_yield_qty
secondary_yield_unit
total_cost
cost_per_yield_unit
cost_per_secondary_unit

---

## 5. `market_list`
Rolls up all ingredients across all recipes.

Columns:

ingredient_id
ingredient_name
unit
cost_per_unit
total_qty
total_cost


