package main

import (
	"log"
	"net"

	"User-Backend/api"
	"User-Backend/internal/services"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("gRPC server listening at %v", listener.Addr())
	}

	grpcServer := grpc.NewServer()

	api.RegisterAuthenticationServiceServer(grpcServer, &services.AuthenticationServer{})
	api.RegisterClassManagerServiceServer(grpcServer, &services.ClassManagerServer{})
	api.RegisterCourseManagerServiceServer(grpcServer, &services.CourseManagerServer{})
	api.RegisterSessionManagerServiceServer(grpcServer, &services.SessionManagerServer{})
	api.RegisterNotificationManagerServiceServer(grpcServer, &services.NotificationManagerServer{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
