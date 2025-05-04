package server

import (
	"User-Backend/api"
	"context"
)

type SessionManagerServer struct {
	api.UnimplementedSessionManagerServiceServer
}

func (s *SessionManagerServer) StartSession(ctx context.Context, req *api.StartSessionRequest) (*api.StartSessionResponse, error) {
	return nil, nil
}

func (s *SessionManagerServer) EndSession(ctx context.Context, req *api.EndSessionRequest) (*api.EndSessionResponse, error) {
	return nil, nil
}

func (s *SessionManagerServer) JoinSession(ctx context.Context, req *api.JoinSessionRequest) (*api.JoinSessionResponse, error) {
	return nil, nil
}
