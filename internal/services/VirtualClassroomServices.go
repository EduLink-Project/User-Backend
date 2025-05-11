package services

import (
	"User-Backend/internal/interfaces"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func CreateVC(data *interfaces.VC, db *sql.DB) error {
	var err error
	if data.Description != "" {
		_, err = db.Exec(`INSERT INTO virtualclassroom (title, description, studentslimit, instructorid, type) VALUES ($1, $2, $3, $4, $5)`, data.Title, data.Description, data.StudentsLimit, data.Instructorid, data.Type)
	} else {
		_, err = db.Exec(`INSERT INTO virtualclassroom (title, studentslimit, instructorid, type) VALUES ($1, $2, $3, $4)`, data.Title, data.StudentsLimit, data.Instructorid, data.Type)
	}
	if err != nil {
		fmt.Printf("Error inserting the recording in the database %v \n", err)
		return errors.New("Error Inserting the record in the database")
	}
	fmt.Println("New Virtual Classroom has been created successfully")
	return nil
}

func GetVCs(role string, id any, db *sql.DB) ([]interfaces.RetrievedVC, error) {
	if role == "student" {
		rows, err := db.Query(`SELECT ei.vcid as id, vc.title as title, vc.description as description, CONCAT(u.first_name, ' ', u.last_name) as instructor, vc.type FROM enrolledin ei join virtualclassroom vc on (ei.vcid = vc.id) join users u on (vc.instructorid = u.id) WHERE ei.studentid = $1`, id)
		if err != nil {
			return nil, errors.New("Error Retrieving Virtual Classrooms")
		}
		defer rows.Close()
		var data []interfaces.RetrievedVC
		for rows.Next() {
			var record interfaces.RetrievedVC
			var description sql.NullString
			err := rows.Scan(&record.ID, &record.Title, &description, &record.Instructor, &record.Type)
			if err != nil {
				return nil, err
			}
			if description.Valid {
				record.Description = description.String
			} else {
				record.Description = ""
			}
			data = append(data, record)
		}
		return data, nil
	} else if role == "instructor" {
		rows, err := db.Query(`SELECT id, title, description, studentsLimit, type FROM virtualclassroom WHERE instructorid = $1`, id)
		if err != nil {
			return nil, errors.New("Error Retrieving Virtual Classrooms")
		}
		defer rows.Close()
		var data []interfaces.RetrievedVC
		for rows.Next() {
			var record interfaces.RetrievedVC
			var description sql.NullString
			err := rows.Scan(&record.ID, &record.Title, &description, &record.StudentsLimit, &record.Type)
			if err != nil {
				return nil, err
			}
			if description.Valid {
				record.Description = description.String
			} else {
				record.Description = ""
			}
			data = append(data, record)
		}
		return data, nil
	} else {
		return nil, errors.New("User's Role is missing")
	}
}

func EnrollStudent(email string, vcid uint64, db *sql.DB) error {
	var studentID int
	err := db.QueryRow("SELECT id from students where email = $1", email).Scan(&studentID)
	if err == sql.ErrNoRows {
		return errors.New("No user found with such email")
	} else if err != nil {
		return errors.New("Error retrieving the student's id")
	}
	_, insertErr := db.Exec("INSERT into enrolledin VALUES ($1, $2)", vcid, studentID)
	if insertErr != nil {
		return errors.New("Error enrolling the student in the Virtual Classroom")
	}
	notification := interfaces.Notifications{
		Type:      "Info",
		Message:   "You have been enrolled in a new Virtual Classroom",
		Title:     "VC Update",
		CreatedAt: time.Now(),
	}
	notificationID, notificationErr := CreateNotification(db, &notification)
	if notificationErr != nil {
		return notificationErr
	}
	notificationErr = InsertNotifications(db, []uint64{uint64(studentID)}, notificationID)
	if notificationErr != nil {
		return notificationErr
	}
	return nil
}
