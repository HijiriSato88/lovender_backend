package service

import (
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
)

type EventsService interface {
	GetUserOshiEvents(userID int64) (*models.OshiEventsResponse, error)
}
type eventsService struct {
	eventsRepo repository.EventsRepository
}

func NewEventsService(eventsRepo repository.EventsRepository) EventsService {
	return &eventsService{eventsRepo: eventsRepo}
}

func (s *eventsService) GetUserOshiEvents(userID int64) (*models.OshiEventsResponse, error) {

	return &models.OshiEventsResponse{
		Oshis: make([]models.OshiEventsResponseItem, 0),
	}, nil
}
