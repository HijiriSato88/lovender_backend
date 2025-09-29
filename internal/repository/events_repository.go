package repository

import (
	"database/sql"
	"fmt"
	"lovender_backend/internal/models"
	"time"
)

type EventsRepository interface {
	GetOshiEventsByUserID(userID int64) (*models.OshiEventsResponse, error)
}

type eventsRepository struct {
	db *sql.DB
}

func NewEventsRepository(db *sql.DB) EventsRepository {
	return &eventsRepository{db: db}
}

func (r *eventsRepository) GetOshiEventsByUserID(userID int64) (*models.OshiEventsResponse, error) {
	query := fmt.Sprintf(`
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
		WHERE o.user_id = %d
		ORDER BY o.id ASC, e.starts_at ASC, e.id ASC
	`, userID)

	rows, err := r.db.Query(query)
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

		// --- 推しバケット確保 ---
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

		// イベントが無い行はスキップ（LEFT JOINの空振り）
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

		// ここまでで idx と ev（models.Event）ができている前提
		dup := false
		for _, existing := range resp.Oshis[idx].Events {
			if existing.ID == ev.ID {
				dup = true
				break
			}
		}
		if !dup {
			resp.Oshis[idx].Events = append(resp.Oshis[idx].Events, ev)
		}

		response.Oshis[index].Events = append(response.Oshis[index].Events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return response, nil
}
