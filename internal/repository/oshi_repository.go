package repository

import (
	"database/sql"
	"lovender_backend/internal/models"
	"sort"
	"time"
)

type OshiRepository interface {
	GetOshisWithDetailsByUserID(userID int64) ([]*models.OshiWithDetails, error)
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
			oshiID, userIDResult                 int64
			accountID, categoryID                *int64
			oshiName, themeColor                 string
			oshiDescription                      *string
			oshiCreatedAt, oshiUpdatedAt         time.Time
			accountURL                           *string
			accountCreatedAt                     *time.Time
			categoryName, categoryDescription    *string
			categoryCreatedAt, categoryUpdatedAt *time.Time
		)

		err := rows.Scan(
			&oshiID, &userIDResult, &oshiName, &oshiDescription, &themeColor,
			&oshiCreatedAt, &oshiUpdatedAt,
			&accountID, &accountURL, &accountCreatedAt,
			&categoryID, &categoryName, &categoryDescription,
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
				ID:        accountID,
				OshiID:    oshiID,
				URL:       accountURL,
				CreatedAt: *accountCreatedAt,
			}
			// 重複チェック
			found := false
			for _, existing := range oshiMap[oshiID].Accounts {
				if existing.ID != nil && accountID != nil && *existing.ID == *accountID {
					found = true
					break
				}
			}
			if !found {
				oshiMap[oshiID].Accounts = append(oshiMap[oshiID].Accounts, account)
			}
		}

		// カテゴリ情報を追加
		if categoryID != nil && categoryName != nil && categoryCreatedAt != nil && categoryUpdatedAt != nil {
			categoryIDUint16 := uint16(*categoryID)
			category := &models.Category{
				ID:          &categoryIDUint16,
				Name:        categoryName,
				Description: categoryDescription,
				CreatedAt:   *categoryCreatedAt,
				UpdatedAt:   *categoryUpdatedAt,
			}
			// 重複チェック
			found := false
			for _, existing := range oshiMap[oshiID].Categories {
				if existing.ID != nil && categoryID != nil && *existing.ID == uint16(*categoryID) {
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
