package internal

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // SQLite driver
)

const DBPath = "db/chefops.db"

// OpenDB opens the SQLite database with foreign keys enabled.
func OpenDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(ON)", DBPath)
	return sql.Open("sqlite", dsn)
}

// LoadRecipeMetadata loads metadata for a recipe by ID
func LoadRecipeMetadata(recipeID int) (*RecipeMetadata, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var metadataStr *string
	query := "SELECT metadata FROM recipes WHERE id = ?"
	err = db.QueryRow(query, recipeID).Scan(&metadataStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("recipe with ID %d not found", recipeID)
		}
		return nil, err
	}

	var metadataStrVal string
	if metadataStr != nil {
		metadataStrVal = *metadataStr
	}

	return LoadMetadata(metadataStrVal)
}

// SaveRecipeMetadata saves metadata for a recipe by ID
func SaveRecipeMetadata(recipeID int, meta *RecipeMetadata) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	metadataStr, err := SaveMetadataToJSON(meta)
	if err != nil {
		return err
	}

	query := "UPDATE recipes SET metadata = ? WHERE id = ?"
	_, err = db.Exec(query, metadataStr, recipeID)
	if err != nil {
		return err
	}

	return nil
}

// GetRecipeIDByName gets recipe ID by name
func GetRecipeIDByName(name string) (int, error) {
	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var id int
	query := "SELECT id FROM recipes WHERE name = ?"
	err = db.QueryRow(query, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("recipe '%s' not found", name)
		}
		return 0, err
	}

	return id, nil
}
