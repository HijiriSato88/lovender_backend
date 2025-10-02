package service

import (
	"errors"
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
)

type OshiGetService interface {
	GetOshiByID(oshiID int64, userID int64) (*models.GetOshiResponse, error)
}

type oshiGetService struct {
	oshiRepo repository.OshiRepository
}

func NewOshiGetService(oshiRepo repository.OshiRepository) OshiGetService {
	return &oshiGetService{
		oshiRepo: oshiRepo,
	}
}

func (s *oshiGetService) GetOshiByID(oshiID int64, userID int64) (*models.GetOshiResponse, error) {
	oshiWithDetails, err := s.oshiRepo.GetOshiByIDAndUserID(oshiID, userID)
	if err != nil {
		if err.Error() == "oshi not found" {
			return nil, errors.New("oshi not found")
		}
		return nil, err
	}

	// URL一覧を配列に変換
	var urls []string
	for _, account := range oshiWithDetails.Accounts {
		urls = append(urls, account.URL)
	}

	// カテゴリ一覧を配列に変換
	var categorySlugs []string
	for _, category := range oshiWithDetails.Categories {
		categorySlugs = append(categorySlugs, category.Slug)
	}

	// レスポンス用の構造体に変換
	resp := &models.GetOshiResponse{
		Oshi: models.GetOshiResponseItem{
			ID:         oshiWithDetails.Oshi.ID,
			Name:       oshiWithDetails.Oshi.Name,
			Color:      oshiWithDetails.Oshi.ThemeColor,
			URLs:       urls,
			Categories: categorySlugs,
		},
	}
	return resp, nil
}
