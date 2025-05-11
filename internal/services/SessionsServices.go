package services

import (
	"User-Backend/internal/interfaces"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func CreateSession(db *sql.DB, data interfaces.Sessions) error {
	datetimeStr := fmt.Sprintf("%s %s", data.DateTime.Date, data.DateTime.Time)
	layout := "2006-01-02 15:04:05"
	dt, err := time.Parse(layout, datetimeStr)
	if err != nil {
		fmt.Printf("Error parsing the DateTime object %v \n", err)
		return errors.New("Error creating the session")
	}
	_, dberr := db.Exec("INSERT INTO sessions(title, dateTime, questionsallowance, vcid, status) VALUES ($1, $2, $3, $4, 'pending')", data.Title, dt, data.QuestionsAllowance, data.Vcid)
	if dberr != nil {
		fmt.Printf("Error executing the query to create the session %v \n", dberr)
		return errors.New("Error creating the session")
	}
	return nil
}

func GetSessions(db *sql.DB, vcid uint64) ([]interfaces.RetrievedSessions, error) {
	rows, err := db.Query("SELECT id, title, status, datetime, questionsallowance from sessions where vcid = $1", vcid)
	if err != nil {
		return nil, errors.New("Error Retrieving Sessions")
	}
	defer rows.Close()
	var data []interfaces.RetrievedSessions
	for rows.Next() {
		var record interfaces.RetrievedSessions
		err := rows.Scan(&record.ID, &record.Title, &record.Status, &record.DateTime, &record.QuestionsAllowance)
		if err != nil {
			return nil, err
		}
		data = append(data, record)
	}
	return data, nil
}

func StartSession(db *sql.DB, ID uint64) error {
	result, err := db.Exec("UPDATE sessions set status = 'live' where id = $1", ID)
	if err != nil {
		fmt.Printf("Error starting the sessions %v \n", err)
		return err
	}
	numOfAffectedRows, execErr := result.RowsAffected()
	if execErr != nil || numOfAffectedRows == 0 {
		fmt.Printf("Error starting the session: Invalid sessions ID\n")
		return errors.New("Error starting the session. No session found with the following ID")
	}
	var studentIDs []uint64
	var vcid uint64
	err = db.QueryRow("SELECT vcid FROM sessions where id = $1", ID).Scan(&vcid)
	if err != nil {
		fmt.Printf("Error starting the sessions %v \n", err)
		return err
	}
	query, queryErr := db.Query("SELECT studentid from enrolledin where vcid = $1", vcid)
	if queryErr != nil {
		fmt.Printf("Error getting the students IDs %v \n", queryErr)
		return queryErr
	}
	for query.Next() {
		var studentID uint64
		queryErr = query.Scan(&studentID)
		if queryErr != nil {
			fmt.Printf("Error getting the students IDs %v \n", queryErr)
			return queryErr
		}
		studentIDs = append(studentIDs, studentID)
	}
	notification := interfaces.Notifications{
		Type:      "Info",
		Message:   "A new session went live!",
		Title:     "VC Update",
		CreatedAt: time.Now(),
	}
	notificationID, notificationErr := CreateNotification(db, &notification)
	if notificationErr != nil {
		return notificationErr
	}
	notificationErr = InsertNotifications(db, studentIDs, notificationID)
	if notificationErr != nil {
		return notificationErr
	}
	return nil
}
