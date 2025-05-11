package server

import (
	"User-Backend/api"
	"User-Backend/internal/interfaces"
	"User-Backend/internal/services"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type SessionManagerServer struct {
	api.UnimplementedSessionManagerServiceServer
	dbCon *sql.DB
}

func NewSessionsManager(db *sql.DB) *SessionManagerServer {
	return &SessionManagerServer{
		dbCon: db,
	}
}

func (s *SessionManagerServer) CreateSession(ctx context.Context, req *api.CreateSessionRequest) (*api.CreateSessionResponse, error) {
	if req.Vcid == 0 || req.Title == "" || req.DateTime == nil || req.DateTime.Date == "" || req.DateTime.Time == "" {
		return nil, errors.New("Failed to create a session. Missing data")
	}
	data := interfaces.Sessions{
		Vcid:               req.Vcid,
		DateTime:           interfaces.DateTime{Date: req.DateTime.Date, Time: req.DateTime.Time},
		QuestionsAllowance: req.QuestionsAllowance,
		Title:              req.Title,
	}
	err := services.CreateSession(s.dbCon, data)
	if err != nil {
		fmt.Printf("Error Create Session : %v \n", err)
		return nil, errors.New("Error Creating Session")
	}
	return &api.CreateSessionResponse{
		Status: "Success",
	}, nil
}

func (s *SessionManagerServer) GetSessions(ctx context.Context, req *api.GetSessionsRequest) (*api.GetSessionsResponse, error) {
	if req.Vcid == nil {
		return nil, errors.New("Error retrieving the sessions. Missing data")
	}
	data, err := services.GetSessions(s.dbCon, req.Vcid.Value)
	if err != nil {
		fmt.Printf("Error getting the sessions %v \n", err)
		return nil, errors.New("Error retrieving the sessions. Internal error")
	}
	var response []*api.Sessions
	for _, value := range data {
		response = append(response, &api.Sessions{
			Id:                 value.ID,
			Title:              value.Title,
			Status:             value.Status,
			DateTime:           timestamppb.New(value.DateTime),
			QuestionsAllowance: value.QuestionsAllowance,
		})
	}
	return &api.GetSessionsResponse{
		Sessions: response,
	}, nil
}

func (s *SessionManagerServer) StartSession(ctx context.Context, req *api.StartSessionRequest) (*api.StartSessionResponse, error) {
	if req.Id == nil {
		return nil, errors.New("Failed to start session. Missing data")
	}
	err := services.StartSession(s.dbCon, req.Id.Value)
	if err != nil {
		return nil, errors.New("Failed to start session. Internal Error")
	}
	return &api.StartSessionResponse{
		Status: "Success",
	}, nil
}
