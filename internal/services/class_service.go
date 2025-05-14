package services

import (
	"User-Backend/api"
	"User-Backend/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

func CreateClass(req *api.CreateClassRequest, db *sql.DB) (*models.Class, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var userExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM instructors WHERE id = $1)", req.UserId).Scan(&userExists)
	if err != nil || !userExists {
		return nil, errors.New("user does not exist or does not have permission to create classes")
	}

	var classID string
	err = tx.QueryRow(`
		INSERT INTO classes(name, user_id, start_time, end_time) 
		VALUES ($1, $2, $3, $4) RETURNING id
	`, req.Name, req.UserId, req.StartTime, req.EndTime).Scan(&classID)

	if err != nil {
		return nil, fmt.Errorf("error creating class: %w", err)
	}

	if len(req.Students) > 0 {

		stmt, err := tx.Prepare("INSERT INTO class_students(class_id, student_id) VALUES ($1, $2)")
		if err != nil {
			return nil, fmt.Errorf("error preparing student enrollment statement: %w", err)
		}
		defer stmt.Close()

		for _, studentID := range req.Students {
			_, err = stmt.Exec(classID, studentID)
			if err != nil {
				return nil, fmt.Errorf("error enrolling student %s: %w", studentID, err)
			}
		}
	}

	if len(req.Files) > 0 {

		stmt, err := tx.Prepare("INSERT INTO class_files(class_id, file_path) VALUES ($1, $2)")
		if err != nil {
			return nil, fmt.Errorf("error preparing class files statement: %w", err)
		}
		defer stmt.Close()

		for _, filePath := range req.Files {
			_, err = stmt.Exec(classID, filePath)
			if err != nil {
				return nil, fmt.Errorf("error adding file %s: %w", filePath, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	class := &models.Class{
		ID:        classID,
		Name:      req.Name,
		UserID:    req.UserId,
		Files:     req.Files,
		Students:  req.Students,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Sessions:  []models.Session{},
	}

	return class, nil
}

func GetClasses(req *api.GetClassesRequest, db *sql.DB) ([]models.Class, error) {

	var isStudent bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM students WHERE id = $1)", req.UserId).Scan(&isStudent)
	if err != nil {
		return nil, fmt.Errorf("error checking user role: %w", err)
	}

	var rows *sql.Rows
	if isStudent {

		rows, err = db.Query(`
			SELECT c.id, c.name, c.user_id, c.start_time, c.end_time
			FROM classes c
			JOIN class_students cs ON c.id = cs.class_id
			WHERE cs.student_id = $1
		`, req.UserId)
	} else {

		rows, err = db.Query(`
			SELECT id, name, user_id, start_time, end_time
			FROM classes
			WHERE user_id = $1
		`, req.UserId)
	}

	if err != nil {
		return nil, fmt.Errorf("error fetching classes: %w", err)
	}
	defer rows.Close()

	classes := []models.Class{}
	for rows.Next() {
		var class models.Class
		err := rows.Scan(&class.ID, &class.Name, &class.UserID, &class.StartTime, &class.EndTime)
		if err != nil {
			return nil, fmt.Errorf("error scanning class row: %w", err)
		}

		filesRows, err := db.Query("SELECT file_path FROM class_files WHERE class_id = $1", class.ID)
		if err == nil {
			defer filesRows.Close()
			files := []string{}
			for filesRows.Next() {
				var file string
				if err := filesRows.Scan(&file); err == nil {
					files = append(files, file)
				}
			}
			class.Files = files
		}

		studentsRows, err := db.Query("SELECT student_id FROM class_students WHERE class_id = $1", class.ID)
		if err == nil {
			defer studentsRows.Close()
			students := []string{}
			for studentsRows.Next() {
				var student string
				if err := studentsRows.Scan(&student); err == nil {
					students = append(students, student)
				}
			}
			class.Students = students
		}

		sessionsRows, err := db.Query("SELECT id, title, date, is_live FROM sessions WHERE class_id = $1", class.ID)
		if err == nil {
			defer sessionsRows.Close()
			sessions := []models.Session{}
			for sessionsRows.Next() {
				var session models.Session
				if err := sessionsRows.Scan(&session.ID, &session.Title, &session.Date, &session.IsLive); err == nil {
					sessions = append(sessions, session)
				}
			}
			class.Sessions = sessions
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func UpdateClass(req *api.UpdateClassRequest, db *sql.DB) (*models.Class, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var isOwner bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM classes WHERE id = $1 AND user_id = $2)",
		req.Classroom.Id, req.UserId).Scan(&isOwner)
	if err != nil || !isOwner {
		return nil, errors.New("user does not have permission to update this class")
	}

	_, err = tx.Exec(`
		UPDATE classes 
		SET name = $1, start_time = $2, end_time = $3
		WHERE id = $4 AND user_id = $5
	`, req.Classroom.Name, req.Classroom.StartTime, req.Classroom.EndTime, req.Classroom.Id, req.UserId)

	if err != nil {
		return nil, fmt.Errorf("error updating class: %w", err)
	}

	_, err = tx.Exec("DELETE FROM class_students WHERE class_id = $1", req.Classroom.Id)
	if err != nil {
		return nil, fmt.Errorf("error removing existing students: %w", err)
	}

	if len(req.Classroom.Students) > 0 {
		stmt, err := tx.Prepare("INSERT INTO class_students(class_id, student_id) VALUES ($1, $2)")
		if err != nil {
			return nil, fmt.Errorf("error preparing student enrollment statement: %w", err)
		}
		defer stmt.Close()

		for _, studentID := range req.Classroom.Students {
			_, err = stmt.Exec(req.Classroom.Id, studentID)
			if err != nil {
				return nil, fmt.Errorf("error enrolling student %s: %w", studentID, err)
			}
		}
	}

	_, err = tx.Exec("DELETE FROM class_files WHERE class_id = $1", req.Classroom.Id)
	if err != nil {
		return nil, fmt.Errorf("error removing existing files: %w", err)
	}

	if len(req.Classroom.Files) > 0 {
		stmt, err := tx.Prepare("INSERT INTO class_files(class_id, file_path) VALUES ($1, $2)")
		if err != nil {
			return nil, fmt.Errorf("error preparing class files statement: %w", err)
		}
		defer stmt.Close()

		for _, filePath := range req.Classroom.Files {
			_, err = stmt.Exec(req.Classroom.Id, filePath)
			if err != nil {
				return nil, fmt.Errorf("error adding file %s: %w", filePath, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	class := &models.Class{}
	class.FromGRPC(req.Classroom)
	return class, nil
}

func DeleteClass(req *api.DeleteClassRequest, db *sql.DB) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var isOwner bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM classes WHERE id = $1 AND user_id = $2)",
		req.ClassroomId, req.UserId).Scan(&isOwner)
	if err != nil || !isOwner {
		return false, errors.New("user does not have permission to delete this class")
	}

	_, err = tx.Exec("DELETE FROM class_students WHERE class_id = $1", req.ClassroomId)
	if err != nil {
		return false, fmt.Errorf("error removing class students: %w", err)
	}

	_, err = tx.Exec("DELETE FROM class_files WHERE class_id = $1", req.ClassroomId)
	if err != nil {
		return false, fmt.Errorf("error removing class files: %w", err)
	}

	_, err = tx.Exec("DELETE FROM sessions WHERE class_id = $1", req.ClassroomId)
	if err != nil {
		return false, fmt.Errorf("error removing class sessions: %w", err)
	}

	_, err = tx.Exec("DELETE FROM classes WHERE id = $1 AND user_id = $2",
		req.ClassroomId, req.UserId)
	if err != nil {
		return false, fmt.Errorf("error deleting class: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
