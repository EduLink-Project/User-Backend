package models

import (
	"User-Backend/api"
	"context"
)

type AuthenticationServer struct {
	api.UnimplementedAuthenticationServiceServer
}

func (s *AuthenticationServer) SignUp(ctx context.Context, req *api.SignUpRequest) (*api.SignUpResponse, error) {
	return nil, nil
}

func (s *AuthenticationServer) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return nil, nil
}

func (s *AuthenticationServer) RefreshToken(ctx context.Context, req *api.RefreshTokenRequest) (*api.RefreshTokenResponse, error) {
	return nil, nil
}

func (s *AuthenticationServer) ValidateToken(ctx context.Context, req *api.ValidateTokenRequest) (*api.ValidateTokenResponse, error) {
	return nil, nil
}

func (s *AuthenticationServer) ForgotPassword(ctx context.Context, req *api.ForgotPasswordRequest) (*api.ForgotPasswordResponse, error) {
	return nil, nil
}