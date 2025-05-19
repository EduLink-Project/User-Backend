package main

import (
	"log"
	"net/http"

	"User-Backend/api"
	"User-Backend/internal/config"
	"User-Backend/internal/handlers"
	"User-Backend/internal/middleware"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

func main() {
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

	// Wrap the gRPC server for both web and native clients
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
		grpcweb.WithAllowedRequestHeaders([]string{"*"}),
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			return true
		}),
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
	)

	// Create HTTP server that handles both gRPC and gRPC-web
	httpServer := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			// Handle CORS preflight requests
			if req.Method == "OPTIONS" {
				resp.Header().Set("Access-Control-Allow-Origin", "*")
				resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				resp.Header().Set("Access-Control-Max-Age", "86400")
				resp.WriteHeader(http.StatusOK)
				return
			}

			// Set CORS headers for all responses
			resp.Header().Set("Access-Control-Allow-Origin", "*")
			resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")

			// Handle both gRPC-web and standard gRPC requests
			if wrappedGrpc.IsGrpcWebRequest(req) || wrappedGrpc.IsAcceptableGrpcCorsRequest(req) {
				wrappedGrpc.ServeHTTP(resp, req)
				return
			}
			resp.WriteHeader(http.StatusNotFound)
		}),
	}

	log.Printf("Starting server on :8080")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
