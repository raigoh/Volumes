package models

import (
	"database/sql"
)

type Category struct {
	ID   int
	Name string
}

func CreateCategory(category *Category) error {
	query := `INSERT INTO categories (name) VALUES (?)`
	result, err := db.Exec(query, category.Name)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	category.ID = int(id)
	return nil
}

func GetCategoryByID(id int) (*Category, error) {
	category := &Category{}
	query := `SELECT category_id, name FROM categories WHERE category_id = ?`
	err := db.QueryRow(query, id).Scan(&category.ID, &category.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return category, err
}

func GetAllCategories() ([]*Category, error) {
	query := `SELECT category_id, name FROM categories`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category := &Category{}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
