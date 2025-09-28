package repository

import (
	"database/sql"
	"lovender_backend/internal/models"
)

type EventsRepository interface {
	GetOshiEventsByUserID(userID int64) (*models.OshiEventsResponse, error)
}

type eventsRepository struct {
	db *sql.DB
}

func NewEventsRepository(db *sql.DB) EventsRepository {
	return &eventsRepository{db: db}
}

func (r *eventsRepository) GetOshiEventsByUserID(userID int64) (*models.OshiEventsResponse, error) {
	return &models.OshiEventsResponse{
		Oshis: make([]models.OshiEventsResponseItem, 0),
	}, nil
}
