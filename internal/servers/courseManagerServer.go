package server

import (
	"User-Backend/api"
	"context"
)

type CourseManagerServer struct {
	api.UnimplementedCourseManagerServiceServer
}

func (s *CourseManagerServer) GetCourses(ctx context.Context, req *api.GetCoursesRequest) (*api.GetCoursesResponse, error) {
	return nil, nil
}
