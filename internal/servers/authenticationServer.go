package server

import (
	"User-Backend/api"
	"User-Backend/internal/services"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type AuthenticationServer struct {
	api.UnimplementedAuthenticationServiceServer
	dbCon *sql.DB
}

func NewAuthServer(db *sql.DB) *AuthenticationServer {
	return &AuthenticationServer{dbCon: db}
}

func (s *AuthenticationServer) SignUp(ctx context.Context, req *api.SignUpRequest) (*api.SignUpResponse, error) {
	if req.Email == "" || req.Firstname == "" || req.Lastname == "" || req.Password == "" || req.Role == "" || req.Username == "" {
		return nil, errors.New("missing data")
	}
	data := map[string]string{
		"password":   req.Password,
		"username":   req.Username,
		"first_name": req.Firstname,
		"last_name":  req.Lastname,
		"email":      req.Email,
		"role":       req.Role,
	}
	jwtToken, err := services.SignUp(data, s.dbCon)
	if err != nil {
		fmt.Println("Error executing the sign up function ", err)
		return nil, errors.New("Error creating a user")
	}
	return &api.SignUpResponse{
		Status: "Success",
		Token:  jwtToken,
	}, nil
}

func (s *AuthenticationServer) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	if req.Password == "" || (req.Email == "" && req.Username == "") || (req.Email != "" && req.Role == "") {
		return nil, errors.New("missing data")
	}
	var data = map[string]string{
		"password": req.Password,
		"username": req.Username,
		"email":    req.Email,
		"role":     req.Role,
	}
	result, err := services.SignIn(data, s.dbCon)
	if err != nil {
		return nil, errors.New("Wrong Credentials")
	}
	return &api.LoginResponse{
		Token: result,
	}, nil
}
