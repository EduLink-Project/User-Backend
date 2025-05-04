package main

import (
	"log"
	"net"

	"User-Backend/api"
	serverConf "User-Backend/internal/config"
	"User-Backend/internal/middleware"
	server "User-Backend/internal/servers"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("gRPC server listening at 3000")
	}

	dbPool := serverConf.InitDB()
	defer dbPool.Close()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor()),
	)

	authServer := server.NewAuthServer(dbPool)

	api.RegisterAuthenticationServiceServer(grpcServer, authServer)
	api.RegisterClassManagerServiceServer(grpcServer, &server.ClassManagerServer{})
	api.RegisterCourseManagerServiceServer(grpcServer, &server.CourseManagerServer{})
	api.RegisterSessionManagerServiceServer(grpcServer, &server.SessionManagerServer{})
	api.RegisterNotificationManagerServiceServer(grpcServer, &server.NotificationManagerServer{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
