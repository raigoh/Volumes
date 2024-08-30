package category

import (
	"fmt"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/pkg/database"
)

// GetPopularCategories retrieves a specified number of categories, ordered by their popularity.
// Popularity is determined by the number of posts associated with each category.
func GetPopularCategories(limit int) ([]models.Category, error) {
	categories := []models.Category{}
	// SQL query to select categories and order them by the count of associated posts
	query := `SELECT c.id, c.name 
              FROM categories c
              JOIN post_categories pc ON c.id = pc.category_id
              GROUP BY c.id
              ORDER BY COUNT(pc.post_id) DESC
              LIMIT ?`

	// Execute the query with the provided limit
	rows, err := database.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the categories slice
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategories retrieves all categories from the database.
func GetCategories() ([]models.Category, error) {
	categories := []models.Category{}
	// Simple query to select all categories
	rows, err := database.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the categories slice
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// InsertInitialCategories populates the database with a predefined list of categories.
// This function is typically used during the initial setup of the application.
func InsertInitialCategories() error {
	// Predefined list of category names
	categories := []string{
		"Fiction",
		"Non-Fiction",
		"Poetry",
		"Science Fiction",
		"Mystery",
		"Fantasy",
		"Romance",
		"Thriller",
		"Horror",
		"Historical Fiction",
		"Biography",
		"Autobiography",
		"Memoir",
		"Self-Help",
		"Philosophy",
		"Psychology",
		"Science",
		"Technology",
		"Business",
		"Economics",
		"Politics",
		"History",
		"Travel",
		"Adventure",
		"Cooking",
		"Art",
		"Music",
		"Drama",
		"Comedy",
		"Crime",
		"Young Adult",
		"Children's Literature",
		"Graphic Novel",
		"Comic Book",
		"Dystopian",
		"Utopian",
		"Alternate History",
		"Steampunk",
		"Cyberpunk",
		"Urban Fantasy",
		"Paranormal",
		"Supernatural",
		"Magical Realism",
		"Literary Fiction",
		"Classic Literature",
		"Contemporary Fiction",
		"Short Story",
		"Novella",
		"Essay",
		"True Crime",
		"Western",
		"Satire",
		"Erotica",
		"Religion",
		"Spirituality",
		"Health",
		"Fitness",
		"Sports",
		"Nature",
		"Environment",
		"Education",
		"Reference",
		"Language",
		"Linguistics",
	}

	// Iterate through the category list and insert each into the database
	for _, category := range categories {
		// Use INSERT OR IGNORE to avoid duplicates if the category already exists
		_, err := database.DB.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", category)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %v", category, err)
		}
	}

	return nil
}
