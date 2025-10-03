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
		SELECT id, slug, name, description, created_at, updated_at
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
		var category models.Category
		var description sql.NullString

		err := rows.Scan(
			&category.ID,
			&category.Slug,
			&category.Name,
			&description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			category.Description = &description.String
		}

		out = append(out, category)
	}
	return out, rows.Err()
}
