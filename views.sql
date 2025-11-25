-- FILE START: views.sql
PRAGMA foreign_keys = ON;

------------------------------------------------------------
-- VIEW 1: RECIPE RAW LINES (ingredients + subrecipes)
------------------------------------------------------------
DROP VIEW IF EXISTS recipe_raw_lines;

CREATE VIEW recipe_raw_lines AS
SELECT
    ri.recipe_id,
    ing.name AS name,
    ri.qty AS qty,
    ing.unit AS unit,
    ing.cost_per_unit AS cost_per_unit,
    ri.qty * ing.cost_per_unit AS line_cost,
    'ingredient' AS type,
    ing.id AS ingredient_id,
    NULL AS subrecipe_id
FROM recipe_items ri
JOIN ingredients ing ON ri.ingredient_id = ing.id

UNION ALL

SELECT
    rs.recipe_id,
    sub.name AS name,
    rs.qty AS qty,
    sub.yield_unit AS unit,
    (
        SELECT SUM(ri2.qty * ing2.cost_per_unit)
        FROM recipe_items ri2
        JOIN ingredients ing2 ON ing2.id = ri2.ingredient_id
        WHERE ri2.recipe_id = sub.id
    ) / sub.yield_qty AS cost_per_unit,
    rs.qty * (
        SELECT SUM(ri2.qty * ing2.cost_per_unit)
        FROM recipe_items ri2
        JOIN ingredients ing2 ON ing2.id = ri2.ingredient_id
        WHERE ri2.recipe_id = sub.id
    ) / sub.yield_qty AS line_cost,
    'subrecipe' AS type,
    NULL AS ingredient_id,
    rs.subrecipe_id AS subrecipe_id
FROM recipe_subrecipes rs
JOIN recipes sub ON rs.subrecipe_id = sub.id;

------------------------------------------------------------
-- VIEW 2: EXPANDED INGREDIENT TOTALS (recursive)
------------------------------------------------------------
DROP VIEW IF EXISTS recipe_items_expanded;

CREATE VIEW recipe_items_expanded AS
WITH RECURSIVE expand(recipe_id, ingredient_id, qty) AS (
    SELECT ri.recipe_id, ri.ingredient_id, ri.qty
    FROM recipe_items ri

    UNION ALL

    SELECT rs.recipe_id, ri2.ingredient_id, rs.qty * ri2.qty
    FROM recipe_subrecipes rs
    JOIN expand ex ON ex.recipe_id = rs.subrecipe_id
    JOIN recipe_items ri2 ON ri2.recipe_id = rs.subrecipe_id
)
SELECT 
    recipe_id, 
    ingredient_id, 
    SUM(qty) AS total_qty
FROM expand
GROUP BY recipe_id, ingredient_id;

------------------------------------------------------------
-- VIEW 2.5: EXPANDED RECIPE ITEMS WITH DETAIL (for export)
------------------------------------------------------------
DROP VIEW IF EXISTS recipe_items_expanded_detail_export;

CREATE VIEW recipe_items_expanded_detail_export AS
SELECT 
    r.name AS recipe_name,
    'ingredient' AS item_type,
    ing.name AS ingredient_name,
    COALESCE(exp.total_qty, ri.qty) AS qty,
    ing.unit AS ingredient_unit,
    COALESCE(exp.total_qty, ri.qty) * ing.cost_per_unit AS line_cost
FROM recipes r
LEFT JOIN recipe_items ri ON r.id = ri.recipe_id
LEFT JOIN ingredients ing ON ri.ingredient_id = ing.id
LEFT JOIN recipe_items_expanded exp ON r.id = exp.recipe_id AND ing.id = exp.ingredient_id

UNION ALL

SELECT 
    r.name AS recipe_name,
    'subrecipe' AS item_type,
    sub.name AS ingredient_name,
    rs.qty AS qty,
    sub.yield_unit AS ingredient_unit,
    rs.qty * (
        SELECT SUM(ri2.qty * ing2.cost_per_unit)
        FROM recipe_items ri2
        JOIN ingredients ing2 ON ing2.id = ri2.ingredient_id
        WHERE ri2.recipe_id = sub.id
    ) / sub.yield_qty AS line_cost
FROM recipes r
LEFT JOIN recipe_subrecipes rs ON r.id = rs.recipe_id
LEFT JOIN recipes sub ON sub.id = rs.subrecipe_id
WHERE sub.id IS NOT NULL;

------------------------------------------------------------
-- VIEW 3: EXPANDED DETAIL VIEW (for scaling & export)
------------------------------------------------------------
DROP VIEW IF EXISTS recipe_items_expanded_detail;

CREATE VIEW recipe_items_expanded_detail AS
SELECT
    r.id AS recipe_id,
    'ingredient' AS item_type,
    ing.name AS ingredient_name,
    ri.qty AS qty,
    ing.unit AS unit
FROM recipe_items ri
JOIN ingredients ing ON ing.id = ri.ingredient_id
JOIN recipes r ON r.id = ri.recipe_id

UNION ALL

SELECT
    r.id AS recipe_id,
    'subrecipe' AS item_type,
    sub.name AS ingredient_name,
    rs.qty AS qty,
    sub.yield_unit AS unit
FROM recipe_subrecipes rs
JOIN recipes r ON r.id = rs.recipe_id
JOIN recipes sub ON sub.id = rs.subrecipe_id;

------------------------------------------------------------
-- VIEW 4: RECIPE TOTAL COSTS
------------------------------------------------------------
DROP VIEW IF EXISTS recipe_totals;

CREATE VIEW recipe_totals AS
SELECT
    r.id AS recipe_id,
    r.name AS recipe_name,
    (
        SELECT SUM(line_cost)
        FROM recipe_raw_lines
        WHERE recipe_id = r.id
    ) AS total_cost,
    r.yield_qty,
    r.yield_unit,
    (
        SELECT SUM(line_cost)
        FROM recipe_raw_lines
        WHERE recipe_id = r.id
    ) / r.yield_qty AS cost_per_yield_unit,
    r.secondary_yield_qty,
    r.secondary_yield_unit,
    CASE WHEN r.secondary_yield_qty > 0 THEN
        (
            SELECT SUM(line_cost)
            FROM recipe_raw_lines
            WHERE recipe_id = r.id
        ) / r.secondary_yield_qty
    END AS cost_per_secondary_unit
FROM recipes r;

------------------------------------------------------------
-- VIEW 5: MARKET LIST (global shopping list)
------------------------------------------------------------
DROP VIEW IF EXISTS market_list;

CREATE VIEW market_list AS
SELECT
    ing.id AS ingredient_id,
    ing.name AS ingredient_name,
    ing.unit AS unit,
    ing.cost_per_unit,
    SUM(exp.total_qty) AS total_qty,
    SUM(exp.total_qty * ing.cost_per_unit) AS total_cost
FROM recipe_items_expanded exp
JOIN ingredients ing ON exp.ingredient_id = ing.id
GROUP BY ing.id, ing.name, ing.unit, ing.cost_per_unit
ORDER BY ing.name;

/* market_list(ingredient_id,ingredient_name,unit,cost_per_unit,total_qty,total_cost) */

-- FILE END: views.sql
