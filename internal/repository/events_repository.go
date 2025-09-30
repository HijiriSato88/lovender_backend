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
