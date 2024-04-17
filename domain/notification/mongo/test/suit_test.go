package test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/mamadeusia/NotificationSrv/config"
	"github.com/mamadeusia/NotificationSrv/domain/notification"
	notifMongo "github.com/mamadeusia/NotificationSrv/domain/notification/mongo"
	"github.com/mamadeusia/NotificationSrv/entity"
	"github.com/stretchr/testify/suite"
)

// integration test suite
type IntTestSuite struct {
	suite.Suite
	notifRepo        *notifMongo.MongoRepository
	ctx              context.Context
	requester1_notif []notification.Notification
	requester2_notif []notification.Notification
	validator1_notif []notification.Notification
	validator2_notif []notification.Notification
}

var (
	AdminId      = 12345
	Requester1Id = 123456
	Requester2Id = 1234567
	Validator1Id = 3456730
	Validator2Id = 595010
)

// this is fully necessary for the suite to run.
func TestIntTestSuite(t *testing.T) {
	suite.Run(t, &IntTestSuite{})
}

func (its *IntTestSuite) SetupSuite() {
	// Load conigurations
	if err := config.Load(); err != nil {
		its.FailNow("MAINROUTINE::LoadConfig:: has failed with error", err)
	}

	its.ctx = context.Background()
	notifRepo, err := notifMongo.NewMongoRepository(its.ctx, config.MongoURL())
	if err != nil {
		its.FailNow("MAINROUTINE::NewMongoRepository:: has failed with error", err)
	}
	its.notifRepo = notifRepo
	its.notifRepo.DropNotificationRepository(its.ctx)

}

func (its *IntTestSuite) TestCreateNotification() {
	its.requester1_notif = []notification.Notification{
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Requester1Id),
				MessageType:   entity.AdminAskQuestion,
				MessageStatus: entity.UnRead,
			},
			MessageDetails: &entity.AdminAskQuestionDetails{
				Question: "send more document ",
			},
		},
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Requester1Id),
				MessageType:   entity.RequestRejected,
				MessageStatus: entity.UnRead,
			},
			MessageDetails: &entity.RequestRejectedDetails{
				Reason: "not enough doc",
				Time:   time.Now().UTC(),
			},
		},
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Requester1Id),
				MessageType:   entity.ValidatorQuestions,
				MessageStatus: entity.Read,
			},
			MessageDetails: &entity.ValidatorQuestionsDetails{
				ValidatorQuestionIds: []string{"1", "2", "3", "4", "5"},
			},
		},
	}

	its.requester2_notif = []notification.Notification{
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Requester2Id),
				MessageType:   entity.AdminAskQuestion,
				MessageStatus: entity.UnRead,
			},
			MessageDetails: &entity.AdminAskQuestionDetails{
				Question: "send more document please",
			},
		},
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Requester2Id),
				MessageType:   entity.RequestRejected,
				MessageStatus: entity.Read,
			},
			MessageDetails: &entity.RequestRejectedDetails{
				Time: time.Now().UTC(),
			},
		},
	}

	its.validator1_notif = []notification.Notification{
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Validator1Id),
				MessageType:   entity.NearRequestFound,
				MessageStatus: entity.Read,
			},
			MessageDetails: &entity.NearRequestFoundDetails{
				RequestID: "3466",
				FullName:  "mobin ghazvini",
			},
		},
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Validator1Id),
				MessageType:   entity.NearRequestFound,
				MessageStatus: entity.UnRead,
			},
			MessageDetails: &entity.NearRequestFoundDetails{
				RequestID: "368",
				FullName:  "mobin ghazvini",
			},
		},
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Validator1Id),
				MessageType:   entity.ElectedAsValidator,
				MessageStatus: entity.Read,
			},
			MessageDetails: &entity.ElectedAsValidatorDetails{
				RequesterFullName: "matin ghazanfari",
			},
		},
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Validator1Id),
				MessageType:   entity.RequesterRespondToValidator,
				MessageStatus: entity.Read,
			},
			MessageDetails: &entity.RequesterRespondToValidatorDetails{
				RequestID: "467",
				StoredMessages: []entity.StoredMessage{
					entity.StoredMessage{
						ChatID:    5799532,
						MessageID: 683,
					},
					entity.StoredMessage{
						ChatID:    44788,
						MessageID: 158,
					},
				},
			},
		},
	}

	its.validator2_notif = []notification.Notification{
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Validator2Id),
				MessageType:   entity.ElectedAsValidator,
				MessageStatus: entity.Read,
			},
			MessageDetails: &entity.ElectedAsValidatorDetails{
				RequesterFullName: "amir lashi",
			},
		},
		notification.Notification{
			BaseMessage: &entity.BaseMessage{
				ID:            "",
				From:          int64(AdminId),
				To:            int64(Validator2Id),
				MessageType:   entity.NearRequestFound,
				MessageStatus: entity.UnRead,
			},
			MessageDetails: &entity.NearRequestFoundDetails{
				RequestID: "467",
				FullName:  "sekine jaki",
			},
		},
	}
	// push requester 1 test data
	for i := range its.requester1_notif {
		notif := its.requester1_notif[i]
		_, err := its.notifRepo.CreateNotification(its.ctx, notif)
		its.Nil(err)
	}
	// push requester 2 test data
	for i := range its.requester2_notif {
		notif := its.requester2_notif[i]
		_, err := its.notifRepo.CreateNotification(its.ctx, notif)
		its.Nil(err)
	}
	// push validator 1 test data
	for i := range its.validator1_notif {
		notif := its.validator1_notif[i]
		_, err := its.notifRepo.CreateNotification(its.ctx, notif)
		its.Nil(err)
	}
	// push validator 2 test data
	for i := range its.validator2_notif {
		notif := its.validator2_notif[i]
		_, err := its.notifRepo.CreateNotification(its.ctx, notif)
		its.Nil(err)
	}

}

func check_notif_equality(recieved []notification.Notification, expected []notification.Notification) bool {
	ret := true
	for i := range recieved {
		r := recieved[i]
		e := expected[i]
		//zeros id
		r.BaseMessage.ID = ""
		e.BaseMessage.ID = ""
		ret = ret && (r.MessageDetails.GetMessageType() == e.MessageDetails.GetMessageType())
		switch r.MessageDetails.GetMessageType() {
		case entity.RequestRejected:
			ret = ret && (reflect.DeepEqual(r.BaseMessage, e.BaseMessage))
			r1, _ := r.MessageDetails.GetMessageDetails()
			e1, _ := e.MessageDetails.GetMessageDetails()
			ret = ret && (r1["reason"] == e1["reason"])
		case entity.RequestApproved:
			ret = ret && (reflect.DeepEqual(r.BaseMessage, e.BaseMessage))
		case entity.AdminAskQuestion:
		case entity.ValidatorQuestions:
		case entity.NearRequestFound:
		case entity.ElectedAsValidator:
		case entity.RequesterRespondToValidator:
			ret = ret && (reflect.DeepEqual(r, e))
		}
	}
	return ret
}

func check_equality(recieved []notification.Notification, expected []notification.Notification, limit int64, offset int64) bool {
	if offset*limit > int64(len(expected)) {
		return check_notif_equality(recieved, expected[offset*limit-limit:])
	} else {
		return check_notif_equality(recieved, expected[offset*limit-limit:offset*limit])
	}
}

func (its *IntTestSuite) TestGetRequesterNotification() {
	notifs, err := its.notifRepo.GetRequesterNotifications(its.ctx, int64(Requester1Id), 5, 1)
	its.Nil(err)
	res := check_equality(notifs, its.requester1_notif, 5, 1)
	its.True(res)

	notifs, err = its.notifRepo.GetRequesterNotifications(its.ctx, int64(Requester1Id), 2, 2)
	its.Nil(err)
	res = check_equality(notifs, its.requester1_notif, 2, 2)
	its.True(res)

	notifs, err = its.notifRepo.GetRequesterNotifications(its.ctx, int64(Requester1Id), 2, 1)
	its.Nil(err)
	res = check_equality(notifs, its.requester1_notif, 2, 1)
	its.True(res)

	notifs, err = its.notifRepo.GetRequesterNotifications(its.ctx, int64(Requester2Id), 5, 1)
	its.Nil(err)
	res = check_equality(notifs, its.requester2_notif, 5, 1)
	its.True(res)

	notifs, err = its.notifRepo.GetRequesterNotifications(its.ctx, int64(Requester2Id), 2, 2)
	its.Nil(err)
	res = check_equality(notifs, its.requester2_notif, 2, 2)
	its.True(res)

	notifs, err = its.notifRepo.GetRequesterNotifications(its.ctx, int64(Requester2Id), 2, 1)
	its.Nil(err)
	res = check_equality(notifs, its.requester2_notif, 2, 1)
	its.True(res)
}

func (its *IntTestSuite) TestGetValidatorNotification() {
	notifs, err := its.notifRepo.GetValidatorNotifications(its.ctx, int64(Validator1Id), 5, 1)
	its.Nil(err)
	res := check_equality(notifs, its.validator1_notif, 5, 1)
	its.True(res)

	notifs, err = its.notifRepo.GetValidatorNotifications(its.ctx, int64(Validator1Id), 2, 2)
	its.Nil(err)
	res = check_equality(notifs, its.validator1_notif, 2, 2)
	its.True(res)

	notifs, err = its.notifRepo.GetValidatorNotifications(its.ctx, int64(Validator1Id), 2, 1)
	its.Nil(err)
	res = check_equality(notifs, its.validator1_notif, 2, 1)
	its.True(res)

	notifs, err = its.notifRepo.GetValidatorNotifications(its.ctx, int64(Validator2Id), 5, 1)
	its.Nil(err)
	res = check_equality(notifs, its.validator2_notif, 5, 1)
	its.True(res)

	notifs, err = its.notifRepo.GetValidatorNotifications(its.ctx, int64(Validator2Id), 2, 2)
	its.Nil(err)
	res = check_equality(notifs, its.validator2_notif, 2, 2)
	its.True(res)

	notifs, err = its.notifRepo.GetValidatorNotifications(its.ctx, int64(Validator2Id), 2, 1)
	its.Nil(err)
	res = check_equality(notifs, its.validator2_notif, 2, 1)
	its.True(res)
}

func calcUnreadNotification(msg []notification.Notification) int64 {
	count := int64(0)
	for i := range msg {
		r := msg[i]
		if r.BaseMessage.MessageStatus == entity.UnRead {
			count = count + 1
		}
	}
	return count
}

func (its *IntTestSuite) TestGetUnreadValidatorNotification() {
	count, err := its.notifRepo.GetUnreadValidatorNotificationCount(its.ctx, int64(Validator1Id))
	its.Nil(err)
	its.Equal(calcUnreadNotification(its.validator1_notif), count)

	count, err = its.notifRepo.GetUnreadValidatorNotificationCount(its.ctx, int64(Validator2Id))
	its.Nil(err)
	its.Equal(calcUnreadNotification(its.validator2_notif), count)

}

func (its *IntTestSuite) TestGetUnreadRequesterNotification() {
	count, err := its.notifRepo.GetUnreadRequesterNotificationCount(its.ctx, int64(Requester1Id))
	its.Nil(err)
	its.Equal(calcUnreadNotification(its.requester1_notif), count)

	count, err = its.notifRepo.GetUnreadRequesterNotificationCount(its.ctx, int64(Requester2Id))
	its.Nil(err)
	its.Equal(calcUnreadNotification(its.requester2_notif), count)

}

func (its *IntTestSuite) TestSetNotificationStatus() {
	notifs, err := its.notifRepo.GetRequesterNotifications(its.ctx, int64(Requester1Id), 5, 1)
	for i := range notifs {
		r := notifs[i]
		if r.BaseMessage.MessageStatus == entity.Read {
			r.BaseMessage.MessageStatus = entity.UnRead
		} else {
			r.BaseMessage.MessageStatus = entity.Read
		}
		its.notifRepo.SetNotificationStatus(its.ctx, r.BaseMessage.ID, r.BaseMessage.MessageStatus)
	}
	its.requester1_notif = notifs
	// now check requster notification
	notifs, err = its.notifRepo.GetRequesterNotifications(its.ctx, int64(Requester1Id), 5, 1)
	its.Nil(err)
	res := check_equality(notifs, its.requester1_notif, 5, 1)
	its.True(res)
}

// req := &pb.CreateRequesterNotificationRequest{
// 	From: 12345,
// 	To:   1234567,
// 	MessageOneof: &pb.CreateRequesterNotificationRequest_RequestRejectedDetails{
// 		RequestRejectedDetails: &pb.RequestRejectedDetails{
// 			Reason: "not enough document",
// 			Time:   timestamppb.New(time.Now()),
// 		},
// 	},
// }

// res := pb.CreateRequesterNotificationResponse{}

// err := its.Nhandler.CreateRequesterNotification(its.ctx, req, &res)

// its.Nil(err)
// its.T().Log(res.Msg)
