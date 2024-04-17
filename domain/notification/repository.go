package notification

import (
	"context"

	"github.com/mamadeusia/NotificationSrv/entity"
)

type Repository interface {
	CreateNotification(ctx context.Context, notif Notification) (*Notification, error)
	CreateBulkNotification(ctx context.Context, notifs []Notification) ([]Notification, error)
	// GetNotifications(ctx context.Context, to int64, limit, offset int64) ([]Notification, error)

	GetRequesterNotifications(ctx context.Context, to int64, limit, offset int64) ([]Notification, error)
	GetValidatorNotifications(ctx context.Context, to int64, limit, offset int64) ([]Notification, error)

	GetUnreadRequesterNotificationCount(ctx context.Context, to int64) (int64, error)
	GetUnreadValidatorNotificationCount(ctx context.Context, to int64) (int64, error)

	// GetUnreadNotificationCount(ctx context.Context, to int64) (int64, error) if the requesterBot and validatorBot meld with each other .

	SetNotificationStatus(ctx context.Context, id string, status entity.MessageStatus) error
}
