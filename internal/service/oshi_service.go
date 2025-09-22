package service

import (
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
)

type OshiService interface {
	GetUserOshis(userID int64) (*models.OshisResponse, error)
}

type oshiService struct {
	oshiRepo repository.OshiRepository
}

func NewOshiService(oshiRepo repository.OshiRepository) OshiService {
	return &oshiService{
		oshiRepo: oshiRepo,
	}
}

func (s *oshiService) GetUserOshis(userID int64) (*models.OshisResponse, error) {
	oshisWithDetails, err := s.oshiRepo.GetOshisWithDetailsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var oshiResponses []models.OshiResponse
	for _, detail := range oshisWithDetails {
		// URL一覧を配列に変換
		var urls []string
		for _, account := range detail.Accounts {
			if account.URL != nil {
				urls = append(urls, *account.URL)
			}
		}

		// カテゴリ一覧を配列に変換
		var categoryNames []string
		for _, category := range detail.Categories {
			if category.Name != nil {
				categoryNames = append(categoryNames, *category.Name)
			}
		}

		// レスポンス用の構造体に変換
		oshiResponse := models.OshiResponse{
			ID:         detail.Oshi.ID,
			Name:       detail.Oshi.Name,
			Color:      detail.Oshi.ThemeColor,
			URLs:       urls,
			Categories: categoryNames,
		}

		oshiResponses = append(oshiResponses, oshiResponse)
	}

	return &models.OshisResponse{
		Oshis: oshiResponses,
	}, nil
}
