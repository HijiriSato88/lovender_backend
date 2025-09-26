package repository

import (
	"database/sql"
	"fmt"
	"log"
	"lovender_backend/internal/models"
	"sort"
	"strings"
	"time"
)

type OshiRepository interface {
	GetOshisWithDetailsByUserID(userID int64) ([]*models.OshiWithDetails, error)
	CreateOshi(oshi *models.Oshi) (int64, error)
	AddAccounts(oshiID int64, urls []string) error
	AddCategories(oshiID int64, categories []string) error
}

type oshiRepository struct {
	db *sql.DB
}

func NewOshiRepository(db *sql.DB) OshiRepository {
	return &oshiRepository{db: db}
}

func (r *oshiRepository) GetOshisWithDetailsByUserID(userID int64) ([]*models.OshiWithDetails, error) {
	query := `
		SELECT 
			o.id as oshi_id,
			o.user_id,
			o.name as oshi_name,
			o.description as oshi_description,
			o.theme_color,
			o.created_at as oshi_created_at,
			o.updated_at as oshi_updated_at,
			oa.id as account_id,
			oa.url as account_url,
			oa.created_at as account_created_at,
			c.id as category_id,
			c.slug as category_slug,
			c.name as category_name,
			c.description as category_description,
			c.created_at as category_created_at,
			c.updated_at as category_updated_at
		FROM oshis o
		LEFT JOIN oshi_accounts oa ON o.id = oa.oshi_id
		LEFT JOIN oshi_categories oc ON o.id = oc.oshi_id
		LEFT JOIN categories c ON oc.category_id = c.id
		WHERE o.user_id = ?
		ORDER BY o.id ASC, oa.created_at ASC, c.name ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	oshiMap := make(map[int64]*models.OshiWithDetails)

	for rows.Next() {
		var (
			oshiID, userIDResult                            int64
			accountID, categoryID                           *int64
			oshiName, themeColor                            string
			oshiDescription                                 *string
			oshiCreatedAt, oshiUpdatedAt                    time.Time
			accountURL                                      *string
			accountCreatedAt                                *time.Time
			categorySlug, categoryName, categoryDescription *string
			categoryCreatedAt, categoryUpdatedAt            *time.Time
		)

		err := rows.Scan(
			&oshiID, &userIDResult, &oshiName, &oshiDescription, &themeColor,
			&oshiCreatedAt, &oshiUpdatedAt,
			&accountID, &accountURL, &accountCreatedAt,
			&categoryID, &categorySlug, &categoryName, &categoryDescription,
			&categoryCreatedAt, &categoryUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 推しがマップに未登録の場合追加
		if _, exists := oshiMap[oshiID]; !exists {
			oshiMap[oshiID] = &models.OshiWithDetails{
				Oshi: &models.Oshi{
					ID:          oshiID,
					UserID:      userIDResult,
					Name:        oshiName,
					Description: oshiDescription,
					ThemeColor:  themeColor,
					CreatedAt:   oshiCreatedAt,
					UpdatedAt:   oshiUpdatedAt,
				},
				Accounts:   []*models.OshiAccount{},
				Categories: []*models.Category{},
			}
		}

		// アカウント情報を追加
		if accountID != nil && accountURL != nil && accountCreatedAt != nil {
			account := &models.OshiAccount{
				ID:        *accountID,
				OshiID:    oshiID,
				URL:       *accountURL,
				CreatedAt: *accountCreatedAt,
			}
			// 重複チェック
			found := false
			for _, existing := range oshiMap[oshiID].Accounts {
				if accountID != nil && existing.ID == *accountID {
					found = true
					break
				}
			}
			if !found {
				oshiMap[oshiID].Accounts = append(oshiMap[oshiID].Accounts, account)
			}
		}

		// カテゴリ情報を追加
		if categoryID != nil && categorySlug != nil && categoryName != nil && categoryCreatedAt != nil && categoryUpdatedAt != nil {
			categoryIDUint16 := uint16(*categoryID)
			category := &models.Category{
				ID:          categoryIDUint16,
				Slug:        *categorySlug,
				Name:        *categoryName,
				Description: categoryDescription,
				CreatedAt:   *categoryCreatedAt,
				UpdatedAt:   *categoryUpdatedAt,
			}
			// 重複チェック
			found := false
			for _, existing := range oshiMap[oshiID].Categories {
				if categoryID != nil && existing.ID == uint16(*categoryID) {
					found = true
					break
				}
			}
			if !found {
				oshiMap[oshiID].Categories = append(oshiMap[oshiID].Categories, category)
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// マップを配列に変換（推しID順でソート）
	var oshiIDs []int64
	for oshiID := range oshiMap {
		oshiIDs = append(oshiIDs, oshiID)
	}
	sort.Slice(oshiIDs, func(i, j int) bool {
		return oshiIDs[i] < oshiIDs[j]
	})

	var result []*models.OshiWithDetails
	for _, oshiID := range oshiIDs {
		result = append(result, oshiMap[oshiID])
	}

	return result, nil
}

// 推しを新規作成
func (r *oshiRepository) CreateOshi(oshi *models.Oshi) (int64, error) {
	query := `
		INSERT INTO oshis (user_id, name, description, theme_color, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()

	result, err := r.db.Exec(query, oshi.UserID, oshi.Name, oshi.Description, oshi.ThemeColor, now, now)
	if err != nil {
		// MySQL エラー内容をログに出す
		log.Printf("CreateOshi ERROR: failed to insert oshi (user_id=%d, name=%s): %v",
			oshi.UserID, oshi.Name, err)
		return 0, fmt.Errorf("failed to insert oshi: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("CreateOshi ERROR: failed to get last insert ID: %v", err)
		return 0, fmt.Errorf("failed to get oshi ID: %w", err)
	}
	return id, nil
}

// 推しのurlを追加
func (r *oshiRepository) AddAccounts(oshiID int64, urls []string) error {
	for _, url := range urls {
		_, err := r.db.Exec(`
          INSERT INTO oshi_accounts (oshi_id, url, created_at) 
          VALUES (?, ?, ?)
      `, oshiID, url, time.Now())
		if err != nil {
			log.Printf("AddAccounts ERROR: failed to insert url=%s for oshi_id=%d: %v",
				url, oshiID, err)
			return fmt.Errorf("failed to insert account url: %w", err)
		}
	}
	return nil
}

// 推しにカテゴリを追加
func (r *oshiRepository) AddCategories(oshiID int64, categories []string) error {
	for _, category := range categories {
		_, err := r.db.Exec(`
            INSERT INTO oshi_categories (oshi_id, category_id) 
            VALUES (?, (SELECT id FROM categories WHERE slug = ?))
        `, oshiID, category)
		if err != nil {
			log.Printf("AddCategories ERROR: failed to insert category=%s for oshi_id=%d: %v",
				category, oshiID, err)

			// 特に category_id=null のケースを検知したい場合
			if strings.Contains(err.Error(), "cannot be null") {
				return fmt.Errorf("invalid category: %s", category)
			}
			return fmt.Errorf("failed to insert category: %w", err)
		}
	}
	return nil
}
