package service

import (
	"fmt"
	"lovender_backend/internal/repository"
	"sync"
)

// KeywordCacheService キーワードキャッシュサービス
type KeywordCacheService struct {
	repository *repository.KeywordRepository
	keywords   []repository.CategoryKeyword
	mu         sync.RWMutex
}

// NewKeywordCacheService コンストラクタ
func NewKeywordCacheService(keywordRepo *repository.KeywordRepository) *KeywordCacheService {
	service := &KeywordCacheService{
		repository: keywordRepo,
		keywords:   make([]repository.CategoryKeyword, 0),
	}

	// 起動時にキーワードをロード
	if err := service.LoadKeywords(); err != nil {
		panic(fmt.Sprintf("Failed to load keywords: %v", err))
	}

	return service
}

// LoadKeywords データベースからキーワードをロード
func (s *KeywordCacheService) LoadKeywords() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	keywords, err := s.repository.GetAllKeywords()
	if err != nil {
		return fmt.Errorf("failed to load keywords from repository: %w", err)
	}

	s.keywords = keywords
	fmt.Printf("Loaded %d keywords into memory\n", len(s.keywords))
	return nil
}

// カテゴリIDのキーワードを取得
func (s *KeywordCacheService) GetKeywordsByCategories(categoryIDs []uint16) []repository.CategoryKeyword {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// カテゴリIDをマップ化（高速検索のため）
	categoryMap := make(map[uint16]bool)
	for _, id := range categoryIDs {
		categoryMap[id] = true
	}

	var result []repository.CategoryKeyword
	for _, keyword := range s.keywords {
		if categoryMap[keyword.CategoryID] {
			result = append(result, keyword)
		}
	}

	return result
}
