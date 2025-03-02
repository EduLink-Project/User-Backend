package models

import (
	"User-Backend/api"
	"context"
)

type NotificationManagerServer struct {
	api.UnimplementedNotificationManagerServiceServer
}

func (s *NotificationManagerServer) GetNotifications(ctx context.Context, req *api.GetNotificationsRequest) (*api.GetNotificationsResponse, error) {
	return nil, nil
}
