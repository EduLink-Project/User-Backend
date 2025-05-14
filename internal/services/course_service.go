package services

import (
	"User-Backend/api"
	"User-Backend/internal/models"
	"database/sql"
	"fmt"
)

func GetCourses(req *api.GetCoursesRequest, db *sql.DB) ([]models.Course, error) {

	var isStudent bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM students WHERE id = $1)", req.UserId).Scan(&isStudent)
	if err != nil {
		return nil, fmt.Errorf("error checking user role: %w", err)
	}

	var rows *sql.Rows
	if isStudent {

		rows, err = db.Query(`
			SELECT c.id, c.name
			FROM courses c
			JOIN course_students cs ON c.id = cs.course_id
			WHERE cs.student_id = $1
		`, req.UserId)
	} else {

		rows, err = db.Query(`
			SELECT id, name
			FROM courses
			WHERE instructor_id = $1
		`, req.UserId)
	}

	if err != nil {
		return nil, fmt.Errorf("error fetching courses: %w", err)
	}
	defer rows.Close()

	courses := []models.Course{}
	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.ID, &course.Name)
		if err != nil {
			return nil, fmt.Errorf("error scanning course row: %w", err)
		}

		sessionsRows, err := db.Query(`
			SELECT id, title, date, is_live
			FROM sessions
			WHERE course_id = $1
		`, course.ID)
		if err == nil {
			defer sessionsRows.Close()
			sessions := []models.Session{}
			for sessionsRows.Next() {
				var session models.Session
				if err := sessionsRows.Scan(&session.ID, &session.Title, &session.Date, &session.IsLive); err == nil {
					sessions = append(sessions, session)
				}
			}
			course.Sessions = sessions
		}

		courses = append(courses, course)
	}

	return courses, nil
}
