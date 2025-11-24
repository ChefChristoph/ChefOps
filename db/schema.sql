PRAGMA foreign_keys = ON;

-- INGREDIENTS:
-- One row per purchasable item from your marketlist / supplier list
CREATE TABLE IF NOT EXISTS ingredients (
    id              INTEGER PRIMARY KEY,
    name            TEXT NOT NULL UNIQUE,
    unit            TEXT NOT NULL,          -- e.g. kg, l, piece
    cost_per_unit   REAL NOT NULL,          -- cost per 1 unit
    notes           TEXT
);

-- RECIPES:
-- One row per recipe (bulk prep or plated dish)
CREATE TABLE IF NOT EXISTS recipes (
    id              INTEGER PRIMARY KEY,
    name            TEXT NOT NULL UNIQUE,
    yield_qty       REAL,                   -- e.g. 10 (kg, l, portions...)
    yield_unit      TEXT,                   -- e.g. kg, l, portion
    notes           TEXT
);

-- RECIPE ITEMS:
-- Links recipes to ingredients with quantities
CREATE TABLE IF NOT EXISTS recipe_items (
    id              INTEGER PRIMARY KEY,
    recipe_id       INTEGER NOT NULL,
    ingredient_id   INTEGER NOT NULL,
    qty             REAL NOT NULL,          -- in ingredient.unit
    notes           TEXT,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS ingredient_conversions (
    id INTEGER PRIMARY KEY,
    ingredient_id INTEGER NOT NULL,
    from_unit TEXT NOT NULL,
    from_qty REAL NOT NULL,
    to_unit TEXT NOT NULL,
    to_qty REAL NOT NULL,
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS recipe_subrecipes (
    id INTEGER PRIMARY KEY,
    recipe_id INTEGER NOT NULL,
    subrecipe_id INTEGER NOT NULL,
    qty REAL NOT NULL,                -- e.g. 3.5 (kg, g, piece, etc.)
    unit TEXT NOT NULL,               -- the unit of the qty
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (subrecipe_id) REFERENCES recipes(id) ON DELETE RESTRICT
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_recipe_items_recipe
    ON recipe_items (recipe_id);

CREATE INDEX IF NOT EXISTS idx_recipe_items_ingredient
    ON recipe_items (ingredient_id);
