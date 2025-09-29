package service

import (
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
)

type EventsService interface {
	GetUserOshiEvents(userID int64) (*models.OshiEventsResponse, error)
	GetEventByID(eventID int64, userID int64) (*models.EventDetailResponse, error)
}
type eventsService struct {
	eventsRepo repository.EventsRepository
}

func NewEventsService(eventsRepo repository.EventsRepository) EventsService {
	return &eventsService{eventsRepo: eventsRepo}
}

func (s *eventsService) GetUserOshiEvents(userID int64) (*models.OshiEventsResponse, error) {
	events, err := s.eventsRepo.GetOshiEventsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *eventsService) GetEventByID(eventID int64, userID int64) (*models.EventDetailResponse, error) {
	// イベント詳細を取得
	eventDetail, err := s.eventsRepo.GetEventByIDWithOshi(eventID, userID)
	if err != nil {
		return nil, err
	}

	return &models.EventDetailResponse{
		Event: *eventDetail,
	}, nil
}
