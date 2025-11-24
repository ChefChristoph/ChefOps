-----------------------------------------------------------------------
-- 1) recipe_raw_lines
--    Ingredients + subrecipes with direct cost resolution
-----------------------------------------------------------------------

DROP VIEW IF EXISTS recipe_raw_lines;

CREATE VIEW recipe_raw_lines AS

-- Ingredient rows
SELECT
    r.id           AS recipe_id,
    r.name         AS recipe_name,
    i.name         AS item_name,
    ri.qty         AS qty,
    i.unit         AS unit,
    (ri.qty * i.cost_per_unit) AS line_cost,
    'ingredient'   AS item_type

FROM recipe_items ri
JOIN recipes r ON r.id = ri.recipe_id
JOIN ingredients i ON i.id = ri.ingredient_id


UNION ALL

-- Subrecipe rows
SELECT
    parent.id       AS recipe_id,
    parent.name     AS recipe_name,
    child.name      AS item_name,
    rs.qty          AS qty,
    rs.unit         AS unit,

    -- subrecipe line cost = qty × (child_total_cost / child_yield)
    rs.qty *
        CASE
            WHEN child.yield_qty IS NOT NULL AND child.yield_qty > 0
            THEN (child_tot.total_cost / child.yield_qty)
            ELSE 0.0
        END AS line_cost,

    'subrecipe' AS item_type

FROM recipe_subrecipes rs
JOIN recipes parent ON parent.id = rs.recipe_id
JOIN recipes child  ON child.id = rs.subrecipe_id

LEFT JOIN (
    SELECT
        r.id AS recipe_id,
        SUM(ri2.qty * i2.cost_per_unit) AS total_cost
    FROM recipes r
    LEFT JOIN recipe_items ri2 ON ri2.recipe_id = r.id
    LEFT JOIN ingredients i2 ON i2.id = ri2.ingredient_id
    GROUP BY r.id
) AS child_tot ON child_tot.recipe_id = child.id;



-----------------------------------------------------------------------
-- 2) recipe_totals
--    Total cost, cost per kg, cost per piece
-----------------------------------------------------------------------

DROP VIEW IF EXISTS recipe_totals;

CREATE VIEW recipe_totals AS
SELECT
    r.id AS recipe_id,
    r.name AS recipe_name,

    r.yield_qty,
    r.yield_unit,
    r.secondary_yield_qty,
    r.secondary_yield_unit,

    COALESCE(SUM(rl.line_cost), 0.0) AS total_cost,

    CASE WHEN r.yield_qty > 0
         THEN COALESCE(SUM(rl.line_cost), 0.0) / r.yield_qty
         ELSE NULL END AS cost_per_yield_unit,

    CASE WHEN r.secondary_yield_qty > 0
         THEN COALESCE(SUM(rl.line_cost), 0.0) / r.secondary_yield_qty
         ELSE NULL END AS cost_per_secondary_unit

FROM recipes r
LEFT JOIN recipe_raw_lines rl ON r.id = rl.recipe_id
GROUP BY r.id;



-----------------------------------------------------------------------
-- 3) recipe_items_expanded
--    For recipe show + export
-----------------------------------------------------------------------

DROP VIEW IF EXISTS recipe_items_expanded;

CREATE VIEW recipe_items_expanded AS
SELECT
    recipe_name,
    item_name AS ingredient_name,
    qty,
    unit AS ingredient_unit,
    line_cost,
    item_type
FROM recipe_raw_lines
ORDER BY recipe_name, item_type, ingredient_name;

