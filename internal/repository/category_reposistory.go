package repository

import (
	"context"
	"database/sql"
	"lovender_backend/internal/models"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]models.Category, error)
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) List(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, slug, name, description FROM categories ORDER BY id`)
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
