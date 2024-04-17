package notification

import (
	"context"

	"github.com/mamadeusia/NotificationSrv/domain/notification"
	"github.com/mamadeusia/NotificationSrv/entity"

	"go-micro.dev/v4/events"
	"go-micro.dev/v4/logger"
)

type EventHandler func(ctx context.Context, e []*events.Event) error

// NotificationService - factory object
type NotificationService struct {
	stream events.Stream
	repo   notification.Repository
}

// NewNotificationService - factory function
func NewNotificationService(rm notification.Repository, stream events.Stream) *NotificationService {
	return &NotificationService{
		repo:   rm,
		stream: stream,
	}
}

// GetMessage - is responsible for business logic and routing to repository
func (n *NotificationService) GetRequesterNotifications(ctx context.Context, req entity.GetRequesterNotificationsRequest) ([]notification.Notification, error) {
	result, err := n.repo.GetRequesterNotifications(ctx, req.To, req.Limit, req.Offset)
	if err != nil {
		logger.Info("SERVICE::GetMessage:: has failed , with error : %+v", err)
		return nil, err
	}

	return result, nil
}

// CreateMessage - is responsible for business logic and routing to repository
func (n *NotificationService) CreateRequesterNotification(ctx context.Context, req entity.CreateRequesterNotificationRequest) error {

	_, err := n.repo.CreateNotification(ctx, notification.Notification{
		BaseMessage: &entity.BaseMessage{
			From:          req.From,
			To:            req.To,
			MessageType:   req.DetailMessage.GetMessageType(),
			MessageStatus: entity.UnRead,
		},
		MessageDetails: req.DetailMessage,
	})
	if err != nil {
		logger.Info("SERVICE::CreateMessage:: has failed, with error : %+v", err)
		return err
	}
	return nil
}

func (n *NotificationService) CreateValidatorNotification(ctx context.Context, req entity.CreateValidatorNotificationRequest) error {

	_, err := n.repo.CreateNotification(ctx, notification.Notification{
		BaseMessage: &entity.BaseMessage{
			From:          req.From,
			To:            req.To,
			MessageType:   req.DetailMessage.GetMessageType(),
			MessageStatus: entity.UnRead,
		},
		MessageDetails: req.DetailMessage,
	})
	if err != nil {
		logger.Info("SERVICE::CreateMessage:: has failed, with error : %+v", err)
		return err
	}
	return nil
}

// CreateBulkMessages - is responsible for business logic and routing to repository
func (n *NotificationService) CreateBulkValidatorNotification(ctx context.Context, req entity.CreateBulkValidatorNotificationRequest) error {
	// msgs := make([]interface{}, len(messages))
	// for i, u := range messages {
	// 	msgs[i] = u
	// }
	var notifs []notification.Notification
	for _, notif := range req.Bulk {
		notifs = append(notifs, notification.Notification{
			BaseMessage: &entity.BaseMessage{
				From:          notif.From,
				To:            notif.To,
				MessageType:   notif.DetailMessage.GetMessageType(),
				MessageStatus: entity.UnRead,
			},
			MessageDetails: notif.DetailMessage,
		})
	}

	// n.repo.CreateBulkNotification(ctx)
	_, err := n.repo.CreateBulkNotification(ctx, notifs)
	if err != nil {
		logger.Info("SERVICE::CreateBulkMessages:: has failed with error : %+v", err)
		return err
	}
	return nil
}

func (n *NotificationService) GetValidatorNotifications(ctx context.Context, req entity.GetValidatorNotificationsRequest) ([]notification.Notification, error) {
	notifs, err := n.repo.GetValidatorNotifications(ctx, req.To, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return notifs, nil
}

func (n *NotificationService) GetRequesterUnreadNotificationsCount(ctx context.Context, req entity.GetRequesterUnreadNotificationsCountRequest) (int64, error) {
	count, err := n.repo.GetUnreadRequesterNotificationCount(ctx, req.To)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (n *NotificationService) GetValidatorUnreadNotificationsCount(ctx context.Context, req entity.GetValidatorUnreadNotificationsCountRequest) (int64, error) {
	count, err := n.repo.GetUnreadValidatorNotificationCount(ctx, req.To)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (n *NotificationService) MarkNotificationStatusAsRead(ctx context.Context, req entity.MarkNotificationStatusAsReadRequest) error {
	if err := n.repo.SetNotificationStatus(ctx, req.ID, entity.Read); err != nil {
		logger.Info("SERVICE::MarkMessagesStatusAsRead, has failed with error : %v", err)
		return err
	}
	return nil
}
