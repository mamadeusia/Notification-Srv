package notificationGrpc

import (
	"context"
	"errors"
	"time"

	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mamadeusia/NotificationSrv/entity"
	"github.com/mamadeusia/NotificationSrv/service/notification"

	pb "github.com/mamadeusia/NotificationSrv/proto"
)

const (
	groupName = "notification"
)

// NotificationHandler - factory object
type Handler struct {
	service *notification.NotificationService
}

func New(srv *notification.NotificationService) *Handler {
	return &Handler{
		service: srv,
	}
}

// CreateRequesterMessage implements NotificationSrv.NotificationSrvHandler
func (n *Handler) CreateRequesterNotification(ctx context.Context, req *pb.CreateRequesterNotificationRequest, res *pb.CreateRequesterNotificationResponse) error {

	createRequesterNotificationRequest := entity.CreateRequesterNotificationRequest{
		From: req.To,
		To:   req.From,
	}

	if req.GetAdminAskQuestionDetails() != nil {
		createRequesterNotificationRequest.DetailMessage = &entity.AdminAskQuestionDetails{
			Question: req.GetAdminAskQuestionDetails().Question,
		}
	} else if req.GetRequestApprovedDetails() != nil {
		createRequesterNotificationRequest.DetailMessage = &entity.RequestApprovedDetails{
			Time: req.GetRequestApprovedDetails().Time.AsTime(),
		}
	} else if req.GetRequestRejectedDetails() != nil {
		createRequesterNotificationRequest.DetailMessage = &entity.RequestRejectedDetails{
			Reason: req.GetRequestRejectedDetails().Reason,
			Time:   req.GetRequestRejectedDetails().Time.AsTime(),
		}
	} else if req.GetValidatorQuestionsDetails() != nil {
		//TODO :: Error handling
		createRequesterNotificationRequest.DetailMessage = &entity.ValidatorQuestionsDetails{
			ValidatorQuestionIds: req.GetValidatorQuestionsDetails().ValidatorQuestionIDs,
		}
	}

	if err := n.service.CreateRequesterNotification(ctx, createRequesterNotificationRequest); err != nil {
		logger.Info("HANDLER:CreateRequesterMessage, has failed with error : %v", err)
	}

	// TODO :have to find an common way for braodcasting failures and successes
	res.Msg = "Success!"
	return nil
}

// CreateValidatorMessage implements NotificationSrv.NotificationSrvHandler
func (n *Handler) CreateValidatorNotification(ctx context.Context, req *pb.CreateValidatorNotificationRequest, res *pb.CreateValidatorNotificationResponse) error {

	createValidatorNotificationRequest := entity.CreateValidatorNotificationRequest{
		From: req.To,
		To:   req.From,
	}

	if req.GetNearRequestFoundDetails() != nil {
		createValidatorNotificationRequest.DetailMessage = &entity.NearRequestFoundDetails{
			RequestID: req.GetNearRequestFoundDetails().GetRequestID(),
			FullName:  req.GetNearRequestFoundDetails().GetFullName(),
		}
	} else if req.GetElectedAsValidatorDetails() != nil {
		createValidatorNotificationRequest.DetailMessage = &entity.ElectedAsValidatorDetails{
			RequesterFullName: req.GetElectedAsValidatorDetails().GetRequesterFullName(),
		}
	} else if req.GetRequesterRespondToValidatorDetails() != nil {
		var storedMessages []entity.StoredMessage
		for _, msg := range req.GetRequesterRespondToValidatorDetails().GetStoredNotifications() {
			storedMessages = append(storedMessages, entity.StoredMessage{
				ChatID:    msg.ChatID,
				MessageID: msg.MessageID,
			})
		}
		createValidatorNotificationRequest.DetailMessage = &entity.RequesterRespondToValidatorDetails{
			RequestID:      req.GetRequesterRespondToValidatorDetails().GetRequestID(),
			StoredMessages: storedMessages,
		}
	}

	if err := n.service.CreateValidatorNotification(ctx, createValidatorNotificationRequest); err != nil {
		logger.Info("HANDLER::CreateValidatorMessage,has failed with error %+v", err)
		return err
	}
	res.Msg = "Success!"
	return nil
}

// CreateValidatorMessage implements NotificationSrv.NotificationSrvHandler
func (n *Handler) CreateBulkValidatorNotification(ctx context.Context, req *pb.CreateBulkValidatorNotificationRequest, res *pb.CreateBulkValidatorNotificationResponse) error {
	var createNotifRequests []entity.CreateValidatorNotificationRequest
	for _, notificationReq := range req.GetCreateValidatorNotificationRequest() {
		createValidatorNotificationRequest := entity.CreateValidatorNotificationRequest{
			From: notificationReq.To,
			To:   notificationReq.From,
		}

		if notificationReq.GetNearRequestFoundDetails() != nil {
			createValidatorNotificationRequest.DetailMessage = &entity.NearRequestFoundDetails{
				RequestID: notificationReq.GetNearRequestFoundDetails().GetRequestID(),
				FullName:  notificationReq.GetNearRequestFoundDetails().GetFullName(),
			}
		} else if notificationReq.GetElectedAsValidatorDetails() != nil {
			createValidatorNotificationRequest.DetailMessage = &entity.ElectedAsValidatorDetails{
				RequesterFullName: notificationReq.GetElectedAsValidatorDetails().GetRequesterFullName(),
			}
		} else if notificationReq.GetRequesterRespondToValidatorDetails() != nil {
			var storedMessages []entity.StoredMessage
			for _, msg := range notificationReq.GetRequesterRespondToValidatorDetails().GetStoredNotifications() {
				storedMessages = append(storedMessages, entity.StoredMessage{
					ChatID:    msg.ChatID,
					MessageID: msg.MessageID,
				})
			}
			createValidatorNotificationRequest.DetailMessage = &entity.RequesterRespondToValidatorDetails{
				RequestID:      notificationReq.GetRequesterRespondToValidatorDetails().GetRequestID(),
				StoredMessages: storedMessages,
			}
		}
		createNotifRequests = append(createNotifRequests, createValidatorNotificationRequest)
	}
	if err := n.service.CreateBulkValidatorNotification(ctx, entity.CreateBulkValidatorNotificationRequest{
		Bulk: createNotifRequests,
	}); err != nil {
		logger.Info("HANDLER::CreateBulkValidatorMessage,has failed with error %+v", err)
		return err
	}

	return nil
}

// GetRequesterMessages implements NotificationSrv.NotificationSrvHandler
func (n *Handler) GetRequesterNotifications(ctx context.Context, req *pb.GetRequesterNotificationsRequest, res *pb.GetRequesterNotificationsResponse) error {

	getRequesterNotificationsRequest := entity.GetRequesterNotificationsRequest{
		To:     req.To,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	result, err := n.service.GetRequesterNotifications(ctx, getRequesterNotificationsRequest)
	if err != nil {
		logger.Info("HANDLER::GetRequesterMessages,has got the error : %v", err)
		return err
	}

	for _, msg := range result {
		var notifPbStatus pb.NotificationStatus
		if msg.BaseMessage.MessageStatus == entity.Read {
			notifPbStatus = pb.NotificationStatus_Read
		} else {
			notifPbStatus = pb.NotificationStatus_UnRead
		}
		notification := &pb.Notification{
			Id:     msg.BaseMessage.ID,
			From:   msg.BaseMessage.From,
			To:     msg.BaseMessage.To,
			Status: notifPbStatus,
		}

		switch msg.BaseMessage.MessageType {
		case entity.AdminAskQuestion:
			msgDetails, err := msg.MessageDetails.GetMessageDetails()
			if err != nil {
				return err
			}
			question, err := GetValueMap[string]("question", msgDetails)
			if err != nil {
				return err
			}
			notification.MessageOneof = &pb.Notification_AdminAskQuestionDetails{
				AdminAskQuestionDetails: &pb.AdminAskQuestionDetails{
					RequestID: "",
					Question:  question,
				},
			}
			// TODO :: change this conversion to generic .
		case entity.RequestRejected:
			msgDetails, err := msg.MessageDetails.GetMessageDetails()
			if err != nil {
				return err
			}

			rejectionReason, err := GetValueMap[string]("reason", msgDetails)
			if err != nil {
				return err
			}
			rejectedTime, err := GetValueMap[time.Time]("time", msgDetails)
			if err != nil {
				return err
			}

			notification.MessageOneof = &pb.Notification_RequestRejectedDetails{
				RequestRejectedDetails: &pb.RequestRejectedDetails{
					Reason: rejectionReason,
					Time:   timestamppb.New(rejectedTime),
				},
			}
		case entity.RequestApproved:
			msgDetails, err := msg.MessageDetails.GetMessageDetails()
			if err != nil {
				return err
			}
			approvedTime, err := GetValueMap[time.Time]("time", msgDetails)
			if err != nil {
				return err
			}
			notification.MessageOneof = &pb.Notification_RequestApprovedDetails{
				RequestApprovedDetails: &pb.RequestApprovedDetails{
					Time: timestamppb.New(approvedTime),
				},
			}
		case entity.ValidatorQuestions:
			msgDetails, err := msg.MessageDetails.GetMessageDetails()
			if err != nil {
				return err
			}
			questionIDs, err := GetValueMap[[]string]("validator_question_ids", msgDetails)
			if err != nil {
				return err
			}
			notification.MessageOneof = &pb.Notification_ValidatorQuestionsDetails{
				ValidatorQuestionsDetails: &pb.ValidatorQuestionsDetails{
					ValidatorQuestionIDs: questionIDs,
				},
			}
		default:
			return errors.New("unknow requester type")

		}

		res.Notifications = append(res.Notifications, notification)
	}
	return nil
}

// GetRequesterMessages implements NotificationSrv.NotificationSrvHandler
func (n *Handler) GetValidatorNotifications(ctx context.Context, req *pb.GetValidatorNotificationsRequest, res *pb.GetValidatorNotificationsResponse) error {

	getValidatorNotificationsRequest := entity.GetValidatorNotificationsRequest{
		To:     req.To,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	result, err := n.service.GetValidatorNotifications(ctx, getValidatorNotificationsRequest)
	if err != nil {
		logger.Info("HANDLER::GetRequesterMessages,has got the error : %v", err)
		return err
	}

	for _, msg := range result {
		var notifPbStatus pb.NotificationStatus
		if msg.BaseMessage.MessageStatus == entity.Read {
			notifPbStatus = pb.NotificationStatus_Read
		} else {
			notifPbStatus = pb.NotificationStatus_UnRead
		}
		notification := &pb.Notification{
			Id:     msg.BaseMessage.ID,
			From:   msg.BaseMessage.From,
			To:     msg.BaseMessage.To,
			Status: notifPbStatus,
		}
		// NearRequestFound
		// ElectedAsValidator
		// RequesterRespondToValidator
		switch msg.BaseMessage.MessageType {
		case entity.NearRequestFound:
			msgDetails, err := msg.MessageDetails.GetMessageDetails()
			if err != nil {
				return err
			}
			requestID, err := GetValueMap[string]("request_id", msgDetails)
			if err != nil {
				return err
			}
			fullName, err := GetValueMap[string]("full_name", msgDetails)
			if err != nil {
				return err
			}
			notification.MessageOneof = &pb.Notification_NearRequestFoundDetails{
				NearRequestFoundDetails: &pb.NearRequestFoundDetails{
					RequestID: requestID,
					FullName:  fullName,
				},
			}
			// TODO :: change this conversion to generic .
		case entity.ElectedAsValidator:
			msgDetails, err := msg.MessageDetails.GetMessageDetails()
			if err != nil {
				return err
			}

			requesterFullName, err := GetValueMap[string]("requester_full_name", msgDetails)
			if err != nil {
				return err
			}

			notification.MessageOneof = &pb.Notification_ElectedAsValidatorDetails{
				ElectedAsValidatorDetails: &pb.ElectedAsValidatorDetails{
					RequesterFullName: requesterFullName,
				},
			}
		case entity.RequesterRespondToValidator:
			msgDetails, err := msg.MessageDetails.GetMessageDetails()
			if err != nil {
				return err
			}
			requestID, err := GetValueMap[string]("request_id", msgDetails)
			if err != nil {
				return err
			}
			storedMessages, err := GetValueMap[[]entity.StoredMessage]("stored_messages", msgDetails)
			if err != nil {
				return err
			}
			var pbStoredNotifications []*pb.StoredNotification
			for _, storedMessage := range storedMessages {
				pbStoredNotifications = append(pbStoredNotifications, &pb.StoredNotification{
					ChatID:    storedMessage.ChatID,
					MessageID: storedMessage.MessageID,
				})
			}

			notification.MessageOneof = &pb.Notification_RequesterRespondToValidatorDetails{
				RequesterRespondToValidatorDetails: &pb.RequesterRespondToValidatorDetails{
					RequestID:           requestID,
					StoredNotifications: pbStoredNotifications,
				},
			}

		default:
			return errors.New("unknow requester type")

		}

		res.Notifications = append(res.Notifications, notification)
	}
	return nil

}

// CreateValidatorQuestions implements NotificationSrv.NotificationSrvHandler
func (n *Handler) GetRequesterUnreadNotificationsCount(ctx context.Context, req *pb.GetRequesterUnreadNotificationsCountRequest, res *pb.GetRequesterUnreadNotificationsCountResponse) error {

	getRequesterUnreadNotificationsCount := entity.GetRequesterUnreadNotificationsCountRequest{
		To: req.To,
	}

	count, err := n.service.GetRequesterUnreadNotificationsCount(ctx, getRequesterUnreadNotificationsCount)
	if err != nil {
		logger.Info("HANDLER::CreateValidatorQuestions, has failed with error : %v", err)
		return err
	}
	res.Count = count

	return nil

}

// CreateValidatorQuestions implements NotificationSrv.NotificationSrvHandler
func (n *Handler) GetValidatorUnreadNotificationsCount(ctx context.Context, req *pb.GetValidatorUnreadNotificationsCountRequest, res *pb.GetValidatorUnreadNotificationsCountResponse) error {

	getValidatorUnreadNotificationsCount := entity.GetValidatorUnreadNotificationsCountRequest{
		To: req.To,
	}

	count, err := n.service.GetValidatorUnreadNotificationsCount(ctx, getValidatorUnreadNotificationsCount)
	if err != nil {
		logger.Info("HANDLER::CreateValidatorQuestions, has failed with error : %v", err)
		return err
	}
	res.Count = count

	return nil
}

func (n *Handler) MarkNotificationStatusAsRead(ctx context.Context, req *pb.MarkNotificationStatusAsReadRequest, res *pb.MarkNotificationStatusAsReadResponse) error {
	markNotificationStatusAsReadRequest := entity.MarkNotificationStatusAsReadRequest{
		ID: req.Id,
	}
	if err := n.service.MarkNotificationStatusAsRead(ctx, markNotificationStatusAsReadRequest); err != nil {
		return err
	}
	res.Msg = "succeed!"
	return nil
}
