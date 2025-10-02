package repository

import (
	"database/sql"
	"fmt"
)

// CategoryKeyword カテゴリキーワード構造体
type CategoryKeyword struct {
	ID         uint64 `db:"id"`
	CategoryID uint16 `db:"category_id"`
	Keyword    string `db:"keyword"`
}

// KeywordRepository キーワードリポジトリ
type KeywordRepository struct {
	db *sql.DB
}

// NewKeywordRepository コンストラクタ
func NewKeywordRepository(db *sql.DB) *KeywordRepository {
	return &KeywordRepository{db: db}
}

// GetAllKeywords 全キーワードを取得
func (r *KeywordRepository) GetAllKeywords() ([]CategoryKeyword, error) {
	query := `
		SELECT id, category_id, keyword
		FROM category_keywords
		ORDER BY category_id, keyword
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query keywords: %w", err)
	}
	defer rows.Close()

	var keywords []CategoryKeyword
	for rows.Next() {
		var keyword CategoryKeyword
		if err := rows.Scan(&keyword.ID, &keyword.CategoryID, &keyword.Keyword); err != nil {
			return nil, fmt.Errorf("failed to scan keyword: %w", err)
		}
		keywords = append(keywords, keyword)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return keywords, nil
}
