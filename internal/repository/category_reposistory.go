package repository

import (
	"database/sql"
	"lovender_backend/internal/models"
)

type CategoryRepository interface {
	GetCategory() ([]models.Category, error)
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetCategory() ([]models.Category, error) {
	query := `
		SELECT id, slug, name, description
		FROM categories
		ORDER BY id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Category
	for rows.Next() {
		var id int16
		var slug, name, desc string
		if err := rows.Scan(&id, &slug, &name, &desc); err != nil {
			return nil, err
		}
		out = append(out, models.Category{
			ID: &id, Slug: &slug, Name: &name, Description: &desc,
		})
	}
	return out, rows.Err()
}
