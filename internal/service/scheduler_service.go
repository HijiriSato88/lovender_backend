package service

import (
	"context"
	"log"
	"time"
)

// 定期実行サービス
type SchedulerService struct {
	eventAutoService *EventAutoService
	ctx              context.Context
	cancel           context.CancelFunc
	ticker           *time.Ticker
}

// コンストラクタ
func NewSchedulerService(eventAutoService *EventAutoService) *SchedulerService {
	ctx, cancel := context.WithCancel(context.Background())

	return &SchedulerService{
		eventAutoService: eventAutoService,
		ctx:              ctx,
		cancel:           cancel,
		ticker:           time.NewTicker(24 * time.Hour),
	}
}

// 定期実行を開始
func (s *SchedulerService) Start() {
	go s.run()
}

// 定期実行のメインループ
func (s *SchedulerService) run() {
	// 起動時に1回実行
	log.Println("Running initial auto event creation on startup")
	s.executeAutoEventCreation()

	// 定期実行ループ
	for {
		select {
		case <-s.ticker.C:
			log.Println("Scheduled auto event creation triggered")
			s.executeAutoEventCreation()
		case <-s.ctx.Done():
			log.Println("Scheduler service stopped")
			return
		}
	}
}

// 自動イベント作成を実行
func (s *SchedulerService) executeAutoEventCreation() {
	// タイムアウト付きのコンテキスト（10分）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	startTime := time.Now()
	log.Printf("Auto event creation started at %s", startTime.Format("2006-01-02 15:04:05"))

	result, err := s.eventAutoService.ProcessAutoEventCreation(ctx)
	if err != nil {
		log.Printf("Auto event creation failed: %v", err)
		return
	}

	duration := time.Since(startTime)
	log.Printf("Auto event creation completed in %v - Processed: %d oshis, Created: %d events, Errors: %d",
		duration, result.ProcessedOshis, result.CreatedEvents, len(result.Errors))

	// エラーがある場合は詳細をログ出力
	if len(result.Errors) > 0 {
		log.Printf("Errors during auto event creation:")
		for i, errMsg := range result.Errors {
			log.Printf("  %d: %s", i+1, errMsg)
		}
	}
}

// 定期実行を停止
func (s *SchedulerService) Stop() {
	log.Println("Stopping scheduler service...")

	if s.ticker != nil {
		s.ticker.Stop()
	}

	if s.cancel != nil {
		s.cancel()
	}

	log.Println("Scheduler service stopped")
}

// 次回実行時刻を取得
func (s *SchedulerService) GetNextRunTime() time.Time {
	return time.Now().Add(24 * time.Hour)
}

// スケジューラーの状態を取得
func (s *SchedulerService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"running":      s.ctx.Err() == nil,
		"interval":     "24 hours",
		"next_run_at":  s.GetNextRunTime().Format("2006-01-02 15:04:05"),
		"current_time": time.Now().Format("2006-01-02 15:04:05"),
	}
}
