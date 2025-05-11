package server

import (
	"User-Backend/api"
	"User-Backend/internal/interfaces"
	"User-Backend/internal/services"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ClassManagerServer struct {
	api.UnimplementedClassManagerServiceServer
	dbCon *sql.DB
}

func NewClassManagerServer(db *sql.DB) *ClassManagerServer {
	return &ClassManagerServer{dbCon: db}
}

func (s *ClassManagerServer) CreateClass(ctx context.Context, req *api.CreateClassRequest) (*api.CreateClassResponse, error) {
	if req.Type == "" || req.Title == "" || req.StudentsLimit == 0 {
		return nil, errors.New("Failed to Create a Virtual Classroom, Missing Data")
	}
	fmt.Printf("Role : %v || ID : %v \n", ctx.Value("role"), ctx.Value("user_id"))
	data := interfaces.VC{
		Title:         req.Title,
		Description:   req.Description,
		StudentsLimit: uint(req.StudentsLimit),
		Instructorid:  ctx.Value("user_id"),
		Type:          req.Type,
	}
	err := services.CreateVC(&data, s.dbCon)
	if err != nil {
		return nil, err
	}
	return &api.CreateClassResponse{
		Message: "Success",
	}, nil
}

func (s *ClassManagerServer) GetStudentClasses(ctx context.Context, req *emptypb.Empty) (*api.GetStudentClassesResponse, error) {
	result, err := services.GetVCs(ctx.Value("role").(string), ctx.Value("user_id"), s.dbCon)
	if err != nil {
		fmt.Printf("Error Retrieving Student's VCs %v", err)
		return nil, errors.New("Failed to retrieve your Virtual Classrooms")
	}
	var response []*api.StudentVC
	for _, value := range result {
		response = append(response, &api.StudentVC{
			Id:          uint64(value.ID),
			Title:       value.Title,
			Description: value.Description,
			Instructor:  value.Instructor,
			Type:        value.Type,
		})
	}
	return &api.GetStudentClassesResponse{
		VCs: response,
	}, nil
}

func (s *ClassManagerServer) GetInstructorClasses(ctx context.Context, req *emptypb.Empty) (*api.GetInstructorClassesResponse, error) {
	result, err := services.GetVCs(ctx.Value("role").(string), ctx.Value("user_id"), s.dbCon)
	if err != nil {
		fmt.Printf("Error Retrieving Instructor's VCs %v", err)
		return nil, errors.New("Failed to retrieve your Virtual Classrooms")
	}
	var response []*api.InstructorVC
	for _, value := range result {
		response = append(response, &api.InstructorVC{
			Id:            uint64(value.ID),
			Title:         value.Title,
			Description:   value.Description,
			StudentsLimit: uint32(value.StudentsLimit),
			Type:          value.Type,
		})
	}
	return &api.GetInstructorClassesResponse{
		VCs: response,
	}, nil
}

func (s *ClassManagerServer) EnrollStudents(ctx context.Context, req *api.EnrollStudentsRequest) (*api.EnrollStudentsReponse, error) {
	if req.StudentEmail == "" || req.Vcid == 0 {
		return nil, errors.New("Can't enroll students. Missing Data")
	}
	err := services.EnrollStudent(req.StudentEmail, req.Vcid, s.dbCon)
	if err != nil {
		fmt.Printf("Error enrolling the student %v \n", err)
		return nil, err
	}
	return &api.EnrollStudentsReponse{
		Status: "Success",
	}, nil
}

// func (s *ClassManagerServer) UpdateClass(ctx context.Context, req *api.UpdateClassRequest) (*api.UpdateClassResponse, error) {
// 	return nil, nil
// }

// func (s *ClassManagerServer) DeleteClass(ctx context.Context, req *api.DeleteClassRequest) (*api.DeleteClassResponse, error) {
// 	return nil, nil
// }
