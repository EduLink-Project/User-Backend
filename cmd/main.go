package main

import (
	"log"
	"net"

	"User-Backend/api"
	"User-Backend/internal/config"
	"User-Backend/internal/handlers"
	"User-Backend/internal/middleware"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("gRPC server listening at 3000")
	}

	dbPool := config.InitDB()
	defer dbPool.Close()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor()),
	)

	authHandler := handlers.NewAuthenticationHandler(dbPool)
	classManagerHandler := handlers.NewClassManagerHandler(dbPool)
	courseManagerHandler := handlers.NewCourseManagerHandler(dbPool)
	sessionManagerHandler := handlers.NewSessionManagerHandler(dbPool)
	notificationManagerHandler := handlers.NewNotificationManagerHandler(dbPool)

	api.RegisterAuthenticationServer(grpcServer, authHandler)
	api.RegisterClassManagerServer(grpcServer, classManagerHandler)
	api.RegisterCourseManagerServer(grpcServer, courseManagerHandler)
	api.RegisterSessionManagerServer(grpcServer, sessionManagerHandler)
	api.RegisterNotificationManagerServer(grpcServer, notificationManagerHandler)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
