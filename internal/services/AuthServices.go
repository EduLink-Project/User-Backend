package services

import (
	Config "User-Backend/internal/config"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int    `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func SignUp(data map[string]string, db *sql.DB) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will rollback if not committed

	var userID int
	err = tx.QueryRow(`
        INSERT INTO users(password, first_name, last_name, username) 
        VALUES ($1, $2, $3, $4) RETURNING id
    `, data["password"], data["first_name"], data["last_name"], data["username"]).Scan(&userID)

	if err != nil {
		fmt.Println("Insert user error:", err)
		return "", err
	}

	// Insert into the appropriate role table
	if data["role"] == "student" {
		_, err = tx.Exec(`INSERT INTO students(id, email) VALUES ($1, $2)`, userID, data["email"])
	} else {
		_, err = tx.Exec(`INSERT INTO instructors(id, email) VALUES ($1, $2)`, userID, data["email"])
	}

	if err != nil {
		fmt.Println("Insert role error:", err)
		return "", err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return "", err
	}

	token, err := generatedJWT_Token(userID, data["role"])
	if err != nil {
		return "", err
	}
	return token, nil
}

func SignIn(data map[string]string, db *sql.DB) (string, error) {
	var result error
	var userID int
	if data["role"] == "student" {
		result = db.QueryRow("SELECT u.id as id from users u join students s on (s.id = u.id) WHERE (s.email = $1 OR u.username = $2) AND u.password = $3", data["email"], data["username"], data["password"]).Scan(&userID)
	} else {
		result = db.QueryRow("SELECT u.id as id from users u join instructors i on (i.id = u.id) WHERE (i.email = $1 OR u.username = $2) AND u.password = $3", data["email"], data["username"], data["password"]).Scan(&userID)
	}
	if result != nil {
		return "", errors.New("No user found with the provided credentials")
	}
	token, tokenErr := generatedJWT_Token(userID, data["role"])
	if tokenErr != nil {
		return "", errors.New("Error signing in")
	}
	return token, nil
}

func generatedJWT_Token(id int, role string) (string, error) {
	claims := CustomClaims{
		UserID: id,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecretKey := Config.GetENVdata("JWT_SECRET")
	if jwtSecretKey == "" {
		return "", errors.New("error generating a token")
	}
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
