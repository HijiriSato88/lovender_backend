package service

import (
	"context"
	"fmt"
	"lovender_backend/internal/repository"
	"sync"
	"time"
)

type KeywordCacheService struct {
	repository  *repository.KeywordRepository
	keywords    []repository.CategoryKeyword
	mu          sync.RWMutex
	lastUpdated time.Time
	ttl         time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewKeywordCacheService(keywordRepo *repository.KeywordRepository) *KeywordCacheService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &KeywordCacheService{
		repository: keywordRepo,
		keywords:   make([]repository.CategoryKeyword, 0),
		ttl:        24 * time.Hour,
		ctx:        ctx,
		cancel:     cancel,
	}

	// 起動時にキーワードをロード
	if err := service.LoadKeywords(); err != nil {
		panic(fmt.Sprintf("Failed to load keywords: %v", err))
	}

	// バックグラウンドでの定期更新を開始
	go service.startPeriodicRefresh()

	return service
}

func (s *KeywordCacheService) LoadKeywords() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	keywords, err := s.repository.GetAllKeywords()
	if err != nil {
		return fmt.Errorf("failed to load keywords from repository: %w", err)
	}

	s.keywords = keywords
	s.lastUpdated = time.Now()
	fmt.Printf("Loaded %d keywords into memory at %s\n", len(s.keywords), s.lastUpdated.Format("2006-01-02 15:04:05"))
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

// バックグラウンドでキャッシュを定期更新
func (s *KeywordCacheService) startPeriodicRefresh() {
	ticker := time.NewTicker(s.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.LoadKeywords(); err != nil {
				fmt.Printf("Failed to refresh keywords cache: %v\n", err)
			}
		case <-s.ctx.Done():
			fmt.Println("Keyword cache refresh goroutine stopped")
			return
		}
	}
}

// Shutdown graceful shutdown
func (s *KeywordCacheService) Shutdown() {
	if s.cancel != nil {
		s.cancel()
	}
}
