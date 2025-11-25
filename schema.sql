-- FILE START: db/schema.sql

PRAGMA foreign_keys = ON;

-- --------------------------
-- INGREDIENTS
-- --------------------------
CREATE TABLE IF NOT EXISTS ingredients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    unit TEXT NOT NULL,
    cost_per_unit REAL NOT NULL,
    notes TEXT
);

-- --------------------------
-- RECIPES
-- --------------------------
CREATE TABLE IF NOT EXISTS recipes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    yield_qty REAL NOT NULL,
    yield_unit TEXT NOT NULL,
    secondary_yield_qty REAL,
    secondary_yield_unit TEXT,
    notes TEXT
);

-- --------------------------
-- RECIPE ITEMS (Ingredients inside Recipes)
-- NO UNIT COLUMN HERE
-- --------------------------
CREATE TABLE IF NOT EXISTS recipe_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recipe_id INTEGER NOT NULL,
    ingredient_id INTEGER NOT NULL,
    qty REAL NOT NULL,
    FOREIGN KEY(recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY(ingredient_id) REFERENCES ingredients(id)
);

-- --------------------------
-- SUBRECIPES (Recipes inside Recipes)
-- --------------------------
CREATE TABLE IF NOT EXISTS recipe_subrecipes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recipe_id INTEGER NOT NULL,
    subrecipe_id INTEGER NOT NULL,
    qty REAL NOT NULL,
    unit TEXT NOT NULL,
    FOREIGN KEY(recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY(subrecipe_id) REFERENCES recipes(id)
);

-- --------------------------
-- OPTIONAL INGREDIENT CONVERSIONS
-- --------------------------
CREATE TABLE IF NOT EXISTS ingredient_conversions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ingredient_id INTEGER NOT NULL,
    from_unit TEXT NOT NULL,
    to_unit TEXT NOT NULL,
    factor REAL NOT NULL,
    FOREIGN KEY(ingredient_id) REFERENCES ingredients(id) ON DELETE CASCADE
);

-- FILE END
