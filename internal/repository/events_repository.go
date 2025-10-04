package repository

import (
	"database/sql"
	"fmt"
	"lovender_backend/internal/models"
	"time"
)

type EventsRepository interface {
	GetOshiEventsByUserID(userID int64) (*models.OshiEventsResponse, error)
	GetEventByIDWithOshi(eventID int64, userID int64) (*models.EventDetail, error)
	UpdateEventByID(eventID int64, userID int64, req *models.UpdateEventData) (*models.UpdatedEventDetail, error)
	CreateEventWithOshi(userID int64, req *models.CreateEventData) (*models.EventDetail, error)
	CheckEventExistsByPostID(postID int64) (bool, error)
	CheckEventExistsByPostIDAndOshiID(postID int64, oshiID int64) (bool, error)
	CreateAutoEvent(oshiID int64, postID int64, title, content string, categoryID *uint16, startsAt time.Time, endsAt *time.Time) error
	GetAllOshisWithAccountsAndCategories() ([]*models.OshiWithDetails, error)
}

type eventsRepository struct {
	db *sql.DB
}

func NewEventsRepository(db *sql.DB) EventsRepository {
	return &eventsRepository{db: db}
}

func (r *eventsRepository) GetOshiEventsByUserID(userID int64) (*models.OshiEventsResponse, error) {
	query := `
		SELECT
			o.id as oshi_id,
			o.name as oshi_name,
			o.theme_color,
			e.id as event_id,
			e.title as event_title,
			e.description as event_description,
			e.url as event_url,
			e.starts_at as event_starts_at,
			e.ends_at as event_ends_at,
			e.has_alarm as event_has_alarm,
			e.notification_timing as event_notification_timing,
			e.has_notification_sent as event_has_notification_sent,
			c.id as category_id,
			c.slug as category_slug,
			c.name as category_name
		FROM oshis o
		LEFT JOIN events e ON o.id = e.oshi_id
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE o.user_id = ?
		ORDER BY o.id ASC, e.starts_at ASC, e.id ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := &models.OshiEventsResponse{
		Oshis: make([]models.OshiEventsResponseItem, 0),
	}
	indexByOshi := make(map[int64]int)

	for rows.Next() {
		var (
			oshiID                   int64
			oshiName                 string
			themeColor               string
			eventID                  *int64
			eventTitle               *string
			eventDescription         *string
			eventURL                 *string
			eventStartsAt            *time.Time
			eventEndsAt              *time.Time
			eventHasAlarm            *bool
			eventNotificationTiming  *string
			eventHasNotificationSent *bool
			categoryID               *int64
			categorySlug             *string
			categoryName             *string
		)

		err := rows.Scan(
			&oshiID, &oshiName, &themeColor,
			&eventID, &eventTitle, &eventDescription, &eventURL, &eventStartsAt, &eventEndsAt,
			&eventHasAlarm, &eventNotificationTiming, &eventHasNotificationSent,
			&categoryID, &categorySlug, &categoryName,
		)

		if err != nil {
			return nil, err
		}

		// 推しごとに結果を集計
		index, exists := indexByOshi[oshiID]
		if !exists {
			response.Oshis = append(response.Oshis, models.OshiEventsResponseItem{
				ID:     oshiID,
				Name:   oshiName,
				Color:  themeColor,
				Events: make([]models.Event, 0),
			})
			index = len(response.Oshis) - 1
			indexByOshi[oshiID] = index
		}

		// イベントが無い行はスキップ
		if eventID == nil {
			continue
		}

		// Category: あれば詰める
		var category *models.CategoryItem
		if categoryID != nil && categorySlug != nil && categoryName != nil {
			category = &models.CategoryItem{
				ID:   *categoryID,
				Slug: *categorySlug,
				Name: *categoryName,
			}
		}

		// モデルへ整形（nilは安全なデフォルトに落とす）
		event := models.Event{
			ID:                    *eventID,
			Title:                 *eventTitle,
			Description:           eventDescription,
			URL:                   eventURL,
			Starts_at:             *eventStartsAt,
			Ends_at:               eventEndsAt,
			Has_alarm:             *eventHasAlarm,
			Notification_timing:   *eventNotificationTiming,
			Has_notification_sent: *eventHasNotificationSent,
			Category:              category,
		}

		// 重複チェック
		found := false
		for _, existing := range response.Oshis[index].Events {
			if existing.ID == event.ID {
				found = true
				break
			}
		}
		if !found {
			response.Oshis[index].Events = append(response.Oshis[index].Events, event)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return response, nil
}

func (r *eventsRepository) GetEventByIDWithOshi(eventID int64, userID int64) (*models.EventDetail, error) {
	query := `
		SELECT
			e.id as event_id,
			e.title as event_title,
			e.description as event_description,
			e.url as event_url,
			e.starts_at as event_starts_at,
			e.ends_at as event_ends_at,
			e.has_alarm as event_has_alarm,
			e.notification_timing as event_notification_timing,
			e.has_notification_sent as event_has_notification_sent,
			o.id as oshi_id,
			o.name as oshi_name,
			o.theme_color as oshi_color
		FROM events e
		INNER JOIN oshis o ON e.oshi_id = o.id
		WHERE e.id = ? AND o.user_id = ?
	`

	row := r.db.QueryRow(query, eventID, userID)

	var (
		eventTitle               string
		eventDescription         *string
		eventURL                 *string
		eventStartsAt            time.Time
		eventEndsAt              *time.Time
		eventHasAlarm            bool
		eventNotificationTiming  string
		eventHasNotificationSent bool
		oshiID                   int64
		oshiName                 string
		oshiColor                string
	)

	err := row.Scan(
		&eventID, &eventTitle, &eventDescription, &eventURL, &eventStartsAt, &eventEndsAt,
		&eventHasAlarm, &eventNotificationTiming, &eventHasNotificationSent,
		&oshiID, &oshiName, &oshiColor,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, err
	}

	// EventOshiを組み立て
	oshi := models.EventOshi{
		ID:    oshiID,
		Name:  oshiName,
		Color: oshiColor,
	}

	// EventDetailを組み立て
	eventDetail := &models.EventDetail{
		ID:                    eventID,
		Title:                 eventTitle,
		Description:           eventDescription,
		URL:                   eventURL,
		Starts_at:             eventStartsAt,
		Ends_at:               eventEndsAt,
		Has_alarm:             eventHasAlarm,
		Notification_timing:   eventNotificationTiming,
		Has_notification_sent: eventHasNotificationSent,
		Oshi:                  oshi,
	}

	return eventDetail, nil
}

func (r *eventsRepository) UpdateEventByID(eventID int64, userID int64, req *models.UpdateEventData) (*models.UpdatedEventDetail, error) {
	// イベントがユーザーの所有する推しか確認
	checkQuery := `
		SELECT EXISTS (
		SELECT 1
		FROM events e
		INNER JOIN oshis o ON e.oshi_id = o.id
		WHERE e.id = ? AND o.user_id = ?
		) AS event_exists
	`

	var count int
	err := r.db.QueryRow(checkQuery, eventID, userID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("event not found")
	}

	// イベント更新
	updateQuery := `
		UPDATE events 
		SET title = ?,
		    description = ?,
		    url = ?,
		    starts_at = ?,
		    ends_at = ?,
		    has_alarm = ?,
		    notification_timing = ?,
		    updated_at = NOW(3)
		WHERE id = ?
	`

	_, err = r.db.Exec(
		updateQuery,
		req.Title,
		req.Description,
		req.URL,
		req.Starts_at,
		req.Ends_at,
		req.Has_alarm,
		req.Notification_timing,
		eventID)
	if err != nil {
		return nil, err
	}

	// 更新されたイベント情報を取得
	selectQuery := `
		SELECT 
			id, title, description, url, starts_at, ends_at, 
			has_alarm, notification_timing, has_notification_sent
		FROM events 
		WHERE id = ?
	`

	row := r.db.QueryRow(selectQuery, eventID)

	var (
		id                  int64
		title               string
		description         *string
		url                 *string
		startsAt            time.Time
		endsAt              *time.Time
		hasAlarm            bool
		notificationTiming  string
		hasNotificationSent bool
	)

	err = row.Scan(&id, &title, &description, &url, &startsAt, &endsAt,
		&hasAlarm, &notificationTiming, &hasNotificationSent)
	if err != nil {
		return nil, err
	}

	return &models.UpdatedEventDetail{
		ID:                    id,
		Title:                 title,
		Description:           description,
		URL:                   url,
		Starts_at:             startsAt,
		Ends_at:               endsAt,
		Has_alarm:             hasAlarm,
		Notification_timing:   notificationTiming,
		Has_notification_sent: hasNotificationSent,
	}, nil
}

func (r *eventsRepository) CreateEventWithOshi(userID int64, req *models.CreateEventData) (*models.EventDetail, error) {
	// 推しがユーザーの所有するものか確認
	checkOshiQuery := `
		SELECT EXISTS (
		SELECT 1
		FROM oshis
		WHERE id = ? AND user_id = ?
		) AS oshi_exists
	`

	var oshiExists int
	err := r.db.QueryRow(checkOshiQuery, req.OshiID, userID).Scan(&oshiExists)
	if err != nil {
		return nil, err
	}
	if oshiExists == 0 {
		return nil, fmt.Errorf("oshi not found")
	}
	// イベント作成
	insertQuery := `
		INSERT INTO events (
			oshi_id, title, description, url,
			starts_at, ends_at, has_alarm, notification_timing
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		insertQuery,
		req.OshiID,
		req.Title,
		req.Description,
		req.URL,
		req.Starts_at,
		req.Ends_at,
		req.Has_alarm,
		req.Notification_timing,
	)
	if err != nil {
		return nil, err
	}

	// 作成されたイベントのIDを取得
	eventID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 作成されたイベント詳細を取得
	return r.GetEventByIDWithOshi(eventID, userID)
}

// 投稿IDでイベントが既に存在するかチェック
func (r *eventsRepository) CheckEventExistsByPostID(postID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM events WHERE post_id = ?) AS event_exists`

	var exists bool
	err := r.db.QueryRow(query, postID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check event existence by post_id: %w", err)
	}

	return exists, nil
}

// 投稿IDと推しIDでイベントが既に存在するかチェック
func (r *eventsRepository) CheckEventExistsByPostIDAndOshiID(postID int64, oshiID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM events WHERE post_id = ? AND oshi_id = ?) AS event_exists`

	var exists bool
	err := r.db.QueryRow(query, postID, oshiID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check event existence by post_id and oshi_id: %w", err)
	}

	return exists, nil
}

// 自動イベント作成
func (r *eventsRepository) CreateAutoEvent(oshiID int64, postID int64, title, content string, categoryID *uint16, startsAt time.Time, endsAt *time.Time) error {
	query := `
		INSERT INTO events (
			oshi_id, category_id, post_id, title, description, 
			starts_at, ends_at, has_alarm, notification_timing
		) VALUES (?, ?, ?, ?, ?, ?, ?, 1, '15m')
	`

	_, err := r.db.Exec(query, oshiID, categoryID, postID, title, content, startsAt, endsAt)
	if err != nil {
		return fmt.Errorf("failed to create auto event: %w", err)
	}

	return nil
}

// 全ユーザーの推し情報を取得（アカウントとカテゴリ付き）
func (r *eventsRepository) GetAllOshisWithAccountsAndCategories() ([]*models.OshiWithDetails, error) {
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
		ORDER BY o.id ASC, oa.created_at ASC, c.name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all oshis: %w", err)
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
			return nil, fmt.Errorf("failed to scan oshi row: %w", err)
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
				if existing.ID == *accountID {
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
				if existing.ID == categoryIDUint16 {
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
		return nil, fmt.Errorf("error iterating oshi rows: %w", err)
	}

	var result []*models.OshiWithDetails
	for _, oshiWithDetails := range oshiMap {
		result = append(result, oshiWithDetails)
	}

	return result, nil
}
