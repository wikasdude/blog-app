package repository

import (
	"database/sql"

	"github.com/wikasdude/blog-backend/model"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) CreateCategory(category *model.Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id` // $1 to avoid sql injection
	return r.DB.QueryRow(query, category.Name).Scan(&category.ID)
}

func (r *CategoryRepository) GetAllCategories() ([]*model.Category, error) {
	query := `SELECT id, name FROM categories`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}
