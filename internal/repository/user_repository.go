package repository

import (
	"database/sql"
	"lovender_backend/internal/models"
)

type UserRepository interface {
	GetByID(id int64) (*models.User, error)
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE id = ?
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(query, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = ?
	`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
