package models

import (
	"User-Backend/api"
	"context"
)

type ClassManagerServer struct {
	api.UnimplementedClassManagerServiceServer
}


func (s *ClassManagerServer) CreateClass(ctx context.Context, req *api.CreateClassRequest) (*api.CreateClassResponse, error) {
	return nil, nil
}

func (s *ClassManagerServer) GetClasses(ctx context.Context, req *api.GetClassesRequest) (*api.GetClassesResponse, error) {
	return nil, nil
}

func (s *ClassManagerServer) UpdateClass(ctx context.Context, req *api.UpdateClassRequest) (*api.UpdateClassResponse, error) {
	return nil, nil
}

func (s *ClassManagerServer) DeleteClass(ctx context.Context, req *api.DeleteClassRequest) (*api.DeleteClassResponse, error) {
	return nil, nil
}