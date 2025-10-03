package cache

import (
	"database/sql"
	"lovender_backend/internal/repository"
	"lovender_backend/internal/service"
)

// キャッシュサービスを管理する構造体
type CacheManager struct {
	KeywordCache *service.KeywordCacheService
}

// キャッシュマネージャーのコンストラクタ
func NewCacheManager(db *sql.DB) *CacheManager {
	// キーワードリポジトリとキャッシュサービスを初期化
	keywordRepo := repository.NewKeywordRepository(db)
	keywordCacheService := service.NewKeywordCacheService(keywordRepo)

	return &CacheManager{
		KeywordCache: keywordCacheService,
	}
}

// 全てのキャッシュサービスを停止
func (cm *CacheManager) Shutdown() {
	if cm.KeywordCache != nil {
		cm.KeywordCache.Shutdown()
	}
}

// キーワードキャッシュサービスを取得
func (cm *CacheManager) GetKeywordCache() *service.KeywordCacheService {
	return cm.KeywordCache
}
