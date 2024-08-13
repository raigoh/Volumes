package category

import (
	"fmt"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/pkg/database"
)

func GetPopularCategories(limit int) ([]models.Category, error) {
	categories := []models.Category{}
	query := `SELECT c.id, c.name 
              FROM categories c
              JOIN post_categories pc ON c.id = pc.category_id
              GROUP BY c.id
              ORDER BY COUNT(pc.post_id) DESC
              LIMIT ?`

	rows, err := database.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func GetCategories() ([]models.Category, error) {
	categories := []models.Category{}
	rows, err := database.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func InsertInitialCategories() error {
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

	for _, category := range categories {
		_, err := database.DB.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", category)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %v", category, err)
		}
	}

	return nil
}
