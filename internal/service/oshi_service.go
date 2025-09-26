package service

import (
	"errors"
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
	"net/http" // ← 追加
	"strings"  // ← 追加
	"time"

	"github.com/labstack/echo/v4"
)

type OshiService interface {
	GetUserOshis(userID int64) (*models.OshisResponse, error)
	CreateOshi(userID int64, req *models.CreateOshiRequest) (*models.CreateOshiResponse, error)
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
			urls = append(urls, account.URL)
		}

		// カテゴリ一覧を配列に変換
		var categoryNames []string
		for _, category := range detail.Categories {
			categoryNames = append(categoryNames, category.Name)
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

// 推しの新規作成
func (s *oshiService) CreateOshi(userID int64, req *models.CreateOshiRequest) (*models.CreateOshiResponse, error) {
	existing, err := s.oshiRepo.GetOshisWithDetailsByUserID(userID)
	if err != nil {
		return nil, err
	}
	for _, o := range existing {
		if o.Oshi.Name == req.Name {
			return nil, errors.New("oshi already exists")
		}
	}

	// 推し情報の作成
	oshi := &models.Oshi{
		UserID:      userID,
		Name:        req.Name,
		Description: nil, // Description is optional and can be set to nil
		ThemeColor:  req.Color,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	oshiID, err := s.oshiRepo.CreateOshi(oshi)
	if err != nil {
		return nil, err
	}

	// 推しのurlの追加
	if len(req.URLs) > 0 {
		err = s.oshiRepo.AddAccounts(oshiID, req.URLs)
		if err != nil {
			return nil, err
		}
	}

	if len(req.Categories) > 0 {
		err = s.oshiRepo.AddCategories(oshiID, req.Categories)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			return nil, err
		}
	}

	// レスポンスの整形
	resp := &models.CreateOshiResponse{
		Oshi: models.CreateOshiResponseItem{
			ID:         oshiID,
			Name:       req.Name,
			Color:      req.Color,
			URLs:       req.URLs,
			Categories: req.Categories,
			CreatedAt:  oshi.CreatedAt,
		},
	}

	return resp, nil
}
