package service

import (
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
)

type CommonService interface {
	GetCommon() (*models.CommonResponse, error)
}

type commonService struct {
	commonRepo repository.CategoryRepository
}

func NewCommonService(commonRepo repository.CategoryRepository) CommonService {
	return &commonService{
		commonRepo: commonRepo,
	}
}

func (s *commonService) GetCommon() (*models.CommonResponse, error) {
	category, err := s.commonRepo.GetCategory()
	if err != nil {
		return nil, err
	}

	return &models.CommonResponse{
		Categories: category,
	}, nil
}
