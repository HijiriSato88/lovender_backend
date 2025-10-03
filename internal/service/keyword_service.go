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
		fmt.Printf("Warning: Failed to load keywords at startup: %v\n", err)
		fmt.Println("Service will continue and attempt to load keywords on first access")
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

	// キャッシュが空の場合、DBから再取得を試行
	if len(s.keywords) == 0 {
		s.mu.RUnlock()
		fmt.Printf("Keywords cache is empty, fetching from database for categories: %v\n", categoryIDs)
		if err := s.LoadKeywords(); err != nil {
			fmt.Printf("Failed to reload keywords from database: %v\n", err)
			return []repository.CategoryKeyword{}
		}
		fmt.Printf("Successfully loaded keywords from database, cache now contains %d keywords\n", len(s.keywords))
		s.mu.RLock()
	} else {
		fmt.Printf("Keywords retrieved from cache for categories: %v (cache contains %d keywords)\n", categoryIDs, len(s.keywords))
	}

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

	fmt.Printf("Found %d matching keywords for categories: %v\n", len(result), categoryIDs)
	return result
}

// バックグラウンドでキャッシュを定期更新
func (s *KeywordCacheService) startPeriodicRefresh() {
	ticker := time.NewTicker(s.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.refreshWithRetry()
		case <-s.ctx.Done():
			fmt.Println("Keyword cache refresh goroutine stopped")
			return
		}
	}
}

// リトライ機能付きキャッシュ更新
func (s *KeywordCacheService) refreshWithRetry() {
	const maxRetries = 3
	const retryDelay = 5 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := s.LoadKeywords(); err != nil {
			fmt.Printf("Failed to refresh keywords cache (attempt %d/%d): %v\n", attempt, maxRetries, err)

			if attempt < maxRetries {
				fmt.Printf("Retrying in %v...\n", retryDelay)
				select {
				case <-time.After(retryDelay):
					continue
				case <-s.ctx.Done():
					fmt.Println("Cache refresh cancelled during retry")
					return
				}
			} else {
				fmt.Printf("All %d attempts failed. Cache will remain unchanged until next refresh cycle.\n", maxRetries)
			}
		} else {
			if attempt > 1 {
				fmt.Printf("Cache refresh succeeded on attempt %d\n", attempt)
			}
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
