package service

import (
	"context"
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
)

type CommonResponse struct {
	Categories []models.Category `json:"categories"`
}

type CommonService interface {
	GetCommon(ctx context.Context) (*CommonResponse, error)
}

type commonService struct {
	repo repository.CategoryRepository
}

func NewCommonService(r repository.CategoryRepository) CommonService {
	return &commonService{repo: r}
}

func (s *commonService) GetCommon(ctx context.Context) (*CommonResponse, error) {
	// まずはダミー
	return &CommonResponse{Categories: []models.Category{}}, nil
}
