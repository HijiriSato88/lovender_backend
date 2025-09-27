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
	CreateOshiWithTransaction(oshi *models.Oshi, urls []string, categories []string) (int64, error)
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

// 推し、アカウント、カテゴリを作成
func (r *oshiRepository) CreateOshiWithTransaction(oshi *models.Oshi, urls []string, categories []string) (int64, error) {
	// トランザクション開始
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("CreateOshiWithTransaction ERROR: failed to begin transaction: %v", err)
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// エラー時ロールバック
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("CreateOshiWithTransaction ERROR: failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	// 推しを作成
	now := time.Now()
	oshiQuery := `
		INSERT INTO oshis (user_id, name, description, theme_color, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := tx.Exec(oshiQuery, oshi.UserID, oshi.Name, oshi.Description, oshi.ThemeColor, now, now)
	if err != nil {
		log.Printf("CreateOshiWithTransaction ERROR: failed to insert oshi (user_id=%d, name=%s): %v",
			oshi.UserID, oshi.Name, err)
		return 0, fmt.Errorf("failed to insert oshi: %w", err)
	}

	oshiID, err := result.LastInsertId()
	if err != nil {
		log.Printf("CreateOshiWithTransaction ERROR: failed to get last insert ID: %v", err)
		return 0, fmt.Errorf("failed to get oshi ID: %w", err)
	}

	// アカウントを一括追加
	if len(urls) > 0 {
		err = r.addAccountsInTransaction(tx, oshiID, urls)
		if err != nil {
			return 0, err
		}
	}

	// カテゴリを一括追加
	if len(categories) > 0 {
		err = r.addCategoriesInTransaction(tx, oshiID, categories)
		if err != nil {
			return 0, err
		}
	}

	// トランザクションをコミット
	if err = tx.Commit(); err != nil {
		log.Printf("CreateOshiWithTransaction ERROR: failed to commit transaction: %v", err)
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("CreateOshiWithTransaction SUCCESS: created oshi_id=%d with %d urls and %d categories",
		oshiID, len(urls), len(categories))
	return oshiID, nil
}

// アカウントを一括追加
func (r *oshiRepository) addAccountsInTransaction(tx *sql.Tx, oshiID int64, urls []string) error {
	if len(urls) == 0 {
		return nil
	}

	// バッチINSERTのためのクエリ構築
	urlCount := len(urls)
	const paramsPerURL = 3 // oshi_id, url, created_at

	valueStrings := make([]string, 0, urlCount)
	valueArgs := make([]interface{}, 0, urlCount*paramsPerURL)
	now := time.Now()

	for _, url := range urls {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, oshiID, url, now)
	}

	query := fmt.Sprintf(`
		INSERT INTO oshi_accounts (oshi_id, url, created_at) 
		VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err := tx.Exec(query, valueArgs...)
	if err != nil {
		log.Printf("addAccountsInTransaction ERROR: failed to insert accounts for oshi_id=%d: %v",
			oshiID, err)
		return fmt.Errorf("failed to insert accounts: %w", err)
	}

	return nil
}

// カテゴリを一括追加
func (r *oshiRepository) addCategoriesInTransaction(tx *sql.Tx, oshiID int64, categories []string) error {
	if len(categories) == 0 {
		return nil
	}

	// カテゴリが存在するかチェック
	placeholders := strings.Repeat("?,", len(categories))
	// 最後のカンマを削除
	placeholders = placeholders[:len(placeholders)-1]

	checkQuery := fmt.Sprintf(`
		SELECT slug FROM categories WHERE slug IN (%s)
	`, placeholders)

	args := make([]interface{}, len(categories))
	for i, category := range categories {
		args[i] = category
	}

	rows, err := tx.Query(checkQuery, args...)
	if err != nil {
		log.Printf("addCategoriesInTransaction ERROR: failed to check categories existence: %v", err)
		return fmt.Errorf("failed to check categories existence: %w", err)
	}
	defer rows.Close()

	// 存在するカテゴリを収集
	existingCategories := make(map[string]bool)
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return fmt.Errorf("failed to scan category slug: %w", err)
		}
		existingCategories[slug] = true
	}

	// 存在しないカテゴリをチェック
	var missingCategories []string
	for _, category := range categories {
		if !existingCategories[category] {
			missingCategories = append(missingCategories, category)
		}
	}

	if len(missingCategories) > 0 {
		return fmt.Errorf("invalid categories: %s", strings.Join(missingCategories, ", "))
	}

	// バッチINSERTのためのクエリ構築
	categoryCount := len(categories)
	const paramsPerCategory = 2 // oshi_id, category (slug)

	valueStrings := make([]string, 0, categoryCount)
	valueArgs := make([]interface{}, 0, categoryCount*paramsPerCategory)

	for _, category := range categories {
		valueStrings = append(valueStrings, "(?, (SELECT id FROM categories WHERE slug = ?))")
		valueArgs = append(valueArgs, oshiID, category)
	}

	query := fmt.Sprintf(`
		INSERT INTO oshi_categories (oshi_id, category_id) 
		VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err = tx.Exec(query, valueArgs...)
	if err != nil {
		log.Printf("addCategoriesInTransaction ERROR: failed to insert categories for oshi_id=%d: %v",
			oshiID, err)
		return fmt.Errorf("failed to insert categories: %w", err)
	}
	return nil
}
