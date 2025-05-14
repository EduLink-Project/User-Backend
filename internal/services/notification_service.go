package services

import (
	"User-Backend/api"
	"User-Backend/internal/models"
	"database/sql"
	"fmt"
)

func GetNotifications(req *api.GetNotificationsRequest, db *sql.DB) ([]models.Notification, error) {

	rows, err := db.Query(`
		SELECT id, title, subtitle, time
		FROM notifications
		WHERE user_id = $1
		ORDER BY time DESC
	`, req.UserId)

	if err != nil {
		return nil, fmt.Errorf("error fetching notifications: %w", err)
	}
	defer rows.Close()

	notifications := []models.Notification{}
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(&notification.ID, &notification.Title, &notification.Subtitle, &notification.Time)
		if err != nil {
			return nil, fmt.Errorf("error scanning notification row: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
