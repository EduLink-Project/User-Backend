package server

import (
	"User-Backend/api"
	"User-Backend/internal/services"
	"context"
	"database/sql"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationManagerServer struct {
	api.UnimplementedNotificationManagerServiceServer
	dbCon *sql.DB
}

func NewNotificationsServer(db *sql.DB) *NotificationManagerServer {
	return &NotificationManagerServer{dbCon: db}
}

func (s *NotificationManagerServer) GetNotifications(ctx context.Context, req *emptypb.Empty) (*api.GetNotificationsResponse, error) {
	result, err := services.GetNotifications(s.dbCon, ctx.Value("user_id").(uint64))
	if err != nil {
		return nil, err
	}
	var response []*api.Notification
	for _, value := range result {
		response = append(response, &api.Notification{
			Id:        value.ID,
			Type:      value.Type,
			Isread:    value.IsRead,
			Message:   value.Message,
			Title:     value.Title,
			CreatedAt: timestamppb.New(value.CreatedAt),
		})
	}
	return &api.GetNotificationsResponse{
		Notifications: response,
	}, nil
}
