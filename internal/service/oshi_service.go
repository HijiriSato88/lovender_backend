package service

import (
	"errors"
	"fmt"
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
	"sort"
	"strings"
	"time"
)

type OshiService interface {
	GetUserOshis(userID int64) (*models.OshisResponse, error)
	CreateOshi(userID int64, req *models.CreateOshiRequest) (*models.CreateOshiResponse, error)
	UpdateOshi(oshiID int64, userID int64, req *models.UpdateOshiRequest) (*models.UpdateOshiResponse, error)
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

	// Service層で推しID順にソート
	sort.Slice(oshisWithDetails, func(i, j int) bool {
		return oshisWithDetails[i].Oshi.ID < oshisWithDetails[j].Oshi.ID
	})

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
	// 推し情報の作成
	// descriptionは未実装のためnil固定. 将来的にreqから受け取るかも
	oshi := &models.Oshi{
		UserID:      userID,
		Name:        req.Name,
		Description: nil,
		ThemeColor:  req.Color,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 推し、アカウント、カテゴリ作成
	oshiID, err := s.oshiRepo.CreateOshiWithTransaction(oshi, req.URLs, req.Categories)
	if err != nil {
		// データベース制約違反をキャッチ
		if isDuplicateKeyError(err) {
			return nil, errors.New("oshi already exists")
		}
		// カテゴリ不正エラーをキャッチ
		if strings.Contains(err.Error(), "invalid categories") {
			return nil, errors.New("invalid categories provided")
		}
		return nil, err
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

// 推しの更新
func (s *oshiService) UpdateOshi(oshiID int64, userID int64, req *models.UpdateOshiRequest) (*models.UpdateOshiResponse, error) {
	// 推し情報の作成
	// descriptionは未実装のためnil固定. 将来的にreqから受け取るかも
	oshi := &models.Oshi{
		UserID:      userID,
		Name:        req.Name,
		Description: nil,
		ThemeColor:  req.Color,
		UpdatedAt:   time.Now(),
	}

	// 推し情報を更新
	err := s.oshiRepo.UpdateOshiWithTransaction(oshiID, userID, oshi, req.URLs, req.Categories)
	if err != nil {
		// 推しが見つからないまたは所有者でない場合
		if strings.Contains(err.Error(), "oshi not found or not owned by user") {
			return nil, errors.New("oshi not found")
		}
		// カテゴリ不正エラーをキャッチ
		if strings.Contains(err.Error(), "invalid categories") {
			return nil, errors.New("invalid categories provided")
		}
		return nil, err
	}

	// 更新後の推し情報を取得
	updatedOshi, err := s.oshiRepo.GetOshiByIDAndUserID(oshiID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated oshi: %w", err)
	}

	// URL一覧を配列に変換
	var urls []string
	for _, account := range updatedOshi.Accounts {
		urls = append(urls, account.URL)
	}

	// カテゴリ一覧を配列に変換
	var categoryNames []string
	for _, category := range updatedOshi.Categories {
		categoryNames = append(categoryNames, category.Name)
	}

	// レスポンスの整形
	resp := &models.UpdateOshiResponse{
		Oshi: models.UpdateOshiResponseItem{
			ID:         updatedOshi.Oshi.ID,
			Name:       updatedOshi.Oshi.Name,
			Color:      updatedOshi.Oshi.ThemeColor,
			URLs:       urls,
			Categories: categoryNames,
			CreatedAt:  updatedOshi.Oshi.CreatedAt,
			UpdatedAt:  updatedOshi.Oshi.UpdatedAt,
		},
	}

	return resp, nil
}

// DB制約違反チェック
func isDuplicateKeyError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "Duplicate entry") ||
		strings.Contains(errMsg, "uq_oshis_user_name") ||
		strings.Contains(errMsg, "UNIQUE constraint failed")
}
