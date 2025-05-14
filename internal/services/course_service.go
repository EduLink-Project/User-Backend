package services

import (
	"User-Backend/api"
	"User-Backend/internal/models"
	"database/sql"
	"fmt"
)

func GetCourses(req *api.GetCoursesRequest, db *sql.DB) ([]models.Course, error) {

	var role string
	err := db.QueryRow("SELECT role FROM users WHERE id = $1", req.UserId).Scan(&role)
	if err != nil {
		return nil, fmt.Errorf("error checking user role: %w", err)
	}

	var rows *sql.Rows
	if role == "student" {

		rows, err = db.Query(`
			SELECT c.id, c.name
			FROM classes c
			JOIN class_students cs ON c.id = cs.class_id
			WHERE cs.student_id = $1
		`, req.UserId)
	} else {

		rows, err = db.Query(`
			SELECT id, name
			FROM classes
			WHERE user_id = $1
		`, req.UserId)
	}

	if err != nil {
		return nil, fmt.Errorf("error fetching classes: %w", err)
	}
	defer rows.Close()

	courses := []models.Course{}
	for rows.Next() {
		var classID, className string
		err := rows.Scan(&classID, &className)
		if err != nil {
			return nil, fmt.Errorf("error scanning class row: %w", err)
		}

		sessionsRows, err := db.Query(`
			SELECT id, title, date, is_live
			FROM sessions
			WHERE class_id = $1
		`, classID)

		if err != nil {
			return nil, fmt.Errorf("error fetching sessions: %w", err)
		}
		defer sessionsRows.Close()

		course := models.Course{
			ID:       classID,
			Name:     className,
			Sessions: []models.Session{},
		}

		for sessionsRows.Next() {
			var session models.Session
			if err := sessionsRows.Scan(&session.ID, &session.Title, &session.Date, &session.IsLive); err == nil {
				course.Sessions = append(course.Sessions, session)
			}
		}

		courses = append(courses, course)
	}

	return courses, nil
}
