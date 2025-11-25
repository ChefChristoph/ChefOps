# Database Schema (SQLite)

This file documents the structure of the SQLite database used by ChefOps.

---

## Tables

### 1) `ingredients`
Stores all purchasable items.

| Column          | Type    | Notes                       |
|-----------------|---------|-----------------------------|
| id              | INTEGER | PK                          |
| name            | TEXT    | Unique                      |
| unit            | TEXT    | kg, liter, piece            |
| cost_per_unit   | REAL    |                              |

---

### 2) `recipes`
Defines all bulk and dish recipes.

| Column               | Type    | Notes                                  |
|----------------------|---------|----------------------------------------|
| id                   | INTEGER | PK                                     |
| name                 | TEXT    | Unique                                 |
| yield_qty            | REAL    | e.g., 1.0                              |
| yield_unit           | TEXT    | e.g., kg                               |
| secondary_yield_qty  | REAL    | Optional                               |
| secondary_yield_unit | TEXT    | Optional                               |
| notes                | TEXT    | Optional description                    |

---

### 3) `recipe_items`
Maps ingredients to recipes.

| Column        | Type    | Notes |
|---------------|---------|-------|
| id            | INTEGER | PK    |
| recipe_id     | INTEGER | FK    |
| ingredient_id | INTEGER | FK    |
| qty           | REAL    | Required |

---

### 4) `recipe_subrecipes`
Defines nested recipes.

| Column        | Type    | Notes                  |
|---------------|---------|------------------------|
| id            | INTEGER | PK                    |
| recipe_id     | INTEGER | Parent recipe          |
| subrecipe_id  | INTEGER | Child recipe           |
| qty           | REAL    | Quantity of subrecipe  |
| unit          | TEXT    | Must match subrecipe   |

---

## Indices
- `ingredients.name` unique  
- `recipes.name` unique  
