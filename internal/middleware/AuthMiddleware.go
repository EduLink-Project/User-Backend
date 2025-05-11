package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"User-Backend/internal/config"

	"github.com/golang-jwt/jwt/v5" // v5 is recommended
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var jwtSecret = config.GetENVdata("JWT_SECRET")

var methodRoleMap = map[string][]string{
	"/userAPI.AuthenticationService/Login":                 {},
	"/userAPI.AuthenticationService/SignUp":                {},
	"/userAPI.ClassManagerService/CreateClass":             {"instructor"},
	"/userAPI.ClassManagerService/GetInstructorClasses":    {"instructor"},
	"/userAPI.ClassManagerService/GetStudentClasses":       {"student"},
	"/userAPI.ClassManagerService/EnrollStudents":          {"instructor"},
	"/userAPI.SessionManagerService/CreateSession":         {"instructor"},
	"/userAPI.SessionManagerService/GetSessions":           {"instructor", "student"},
	"/userAPI.SessionManagerService/StartSession":          {"instructor"},
	"/userAPI.NotificationManagerService/GetNotifications": {"instructor", "student"},
}

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		allowedRoles, methodProtected := methodRoleMap[info.FullMethod]
		// If method is not listed, block by default
		if !methodProtected {
			fmt.Printf("Method %s is not registered\n", info.FullMethod)
			return nil, fmt.Errorf("Unauthorized")
		}

		// If allowedRoles is empty → public route (no auth required)
		if len(allowedRoles) == 0 {
			return handler(ctx, req)
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			fmt.Println("Missing Metadata")
			return nil, errors.New("missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			fmt.Println("Missing Authorization Token")
			return nil, errors.New("Not Authenticated")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			fmt.Println("Invalid Token")
			return nil, errors.New("Not Authenticated")
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("Missing Token Payload Data")
			return nil, errors.New("Not Authenticated")
		}
		role, RoleOk := claims["role"].(string)
		user_id, user_idOk := claims["sub"].(float64)
		if !RoleOk {
			fmt.Println("Missing Role in the Token Payload Data")
			return nil, errors.New("Not Authenticated")
		}
		if !user_idOk {
			fmt.Println("Missing ID in the Token Payload Data")
			return nil, errors.New("Not Authenticated")
		}

		// Check if user's role is allowed
		for _, allowed := range allowedRoles {
			if role == allowed {
				ctx = context.WithValue(ctx, "user_id", uint64(user_id))
				ctx = context.WithValue(ctx, "role", role)
				return handler(ctx, req)
			}
		}
		fmt.Printf("unauthorized: role '%s' cannot access %s\n", role, info.FullMethod)
		return nil, fmt.Errorf("Unauthorized")
	}
}
