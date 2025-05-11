package services

import (
	"User-Backend/internal/interfaces"
	"database/sql"
	"errors"
	"fmt"
)

func GetNotifications(db *sql.DB, userid uint64) ([]interfaces.Notifications, error) {
	rows, err := db.Query("SELECT n.*, nb.isread from notifiedBy nb join notifications n on (n.id = nb.notificationid) where nb.userid = $1", userid)
	defer rows.Close()
	if err != nil {
		fmt.Printf("Error Retrieving Notifications : %v \n", err)
		return nil, errors.New("Error returning notifications")
	}
	var data []interfaces.Notifications
	for rows.Next() {
		var temp interfaces.Notifications
		scanErr := rows.Scan(&temp.ID, &temp.Type, &temp.Message, &temp.Title, &temp.CreatedAt, &temp.IsRead)
		if scanErr != nil {
			fmt.Printf("Error Retrieving Notifications : %v \n", scanErr)
			return nil, errors.New("Error returning notifications")
		}
		data = append(data, temp)
	}
	return data, nil
}

func CreateNotification(db *sql.DB, notification *interfaces.Notifications) (uint64, error) {
	var notificationID int64
	err := db.
		QueryRow(
			`INSERT INTO notifications(notificationtype, message, title, createdat)
             VALUES ($1, $2, $3, $4)
             RETURNING id`,
			notification.Type,
			notification.Message,
			notification.Title,
			notification.CreatedAt,
		).
		Scan(&notificationID)
	if err != nil {
		return 0, fmt.Errorf("error creating notification: %w", err)
	}
	return uint64(notificationID), nil
}

func InsertNotifications(db *sql.DB, userIDs []uint64, notificationID uint64) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("Error starting transaction: %v\n", err)
		return err
	}

	// rollback only if commit doesn't happen
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			tx.Rollback()
		}
	}()

	for _, userID := range userIDs {
		_, err = tx.Exec("INSERT INTO notifiedby VALUES($1, $2, $3)", userID, notificationID, false)
		if err != nil {
			fmt.Printf("Error inserting into notifiedby: %v\n", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		return err
	}

	return nil
}
