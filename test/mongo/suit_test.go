package test

import (
	"context"
	"testing"
	"time"

	"github.com/mamadeusia/NotificationSrv/config"
	notifMongo "github.com/mamadeusia/NotificationSrv/domain/notification/mongo"
	questionMongo "github.com/mamadeusia/NotificationSrv/domain/question/mongo"
	"github.com/mamadeusia/NotificationSrv/handler/notificationEvent"
	"github.com/mamadeusia/NotificationSrv/handler/notificationGrpc"
	questionGrpc "github.com/mamadeusia/NotificationSrv/handler/questionGrpc"
	pb "github.com/mamadeusia/NotificationSrv/proto"
	notifService "github.com/mamadeusia/NotificationSrv/service/notification"
	questionService "github.com/mamadeusia/NotificationSrv/service/question"
	"github.com/mamadeusia/go-micro-plugins/events/natsjs"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// integration test suite
type IntTestSuite struct {
	suite.Suite
	notifRepo    *notifMongo.MongoRepository
	questionRepo *questionMongo.MongoRepository
	Nhandler     *notificationGrpc.Handler
	Qhandler     *questionGrpc.Handler
	ctx          context.Context
}

var (
	service = "notificationsrv"
	version = "latest"
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

	questionRepo, err := questionMongo.NewMongoRepository(its.ctx, config.MongoURL())
	if err != nil {
		its.FailNow("MAINROUTINE::NewMongoRepository:: has failed with error ", err)
	}

	its.questionRepo = questionRepo

	stream, err := natsjs.NewStream(
		natsjs.Address(config.NatsURL()),
	)
	if err != nil {
		its.FailNow("natjs stream creation failed", err)
	}

	store, err := natsjs.NewStore(
		natsjs.Address(config.NatsURL()),
	)

	if err != nil {
		its.FailNow("natjs store creation failed", err)
	}

	notificationService := notifService.NewNotificationService(notifRepo, stream)
	quesitonService := questionService.NewQuestionService(questionRepo, stream)

	its.Nhandler = notificationGrpc.New(notificationService)
	its.Qhandler = questionGrpc.New(quesitonService)
	EventHandler, err := notificationEvent.New(config.NatsNotificationPrefix(), config.NatsNotificationRetryDurations(), store, stream)
	if err != nil {
		its.FailNow("notification event creation failed", err)
	}

	if err := EventHandler.StartConsumer(its.ctx); err != nil {
		its.FailNow("unable to start Consumer", err)
	}

}

func (its *IntTestSuite) TestCreateRequsterAdminAskQuestion() {
	req := &pb.CreateRequesterNotificationRequest{
		From: 12345,
		To:   1234567,
		MessageOneof: &pb.CreateRequesterNotificationRequest_AdminAskQuestionDetails{
			AdminAskQuestionDetails: &pb.AdminAskQuestionDetails{
				RequestID: "AABBCCDD",
				Question:  "please add more pic and document",
			},
		},
	}

	res := pb.CreateRequesterNotificationResponse{}

	err := its.Nhandler.CreateRequesterNotification(its.ctx, req, &res)

	its.Nil(err)
	its.T().Log(res.Msg)

}

func (its *IntTestSuite) TestCreateRequesterRequestRejected() {
	req := &pb.CreateRequesterNotificationRequest{
		From: 12345,
		To:   1234567,
		MessageOneof: &pb.CreateRequesterNotificationRequest_RequestRejectedDetails{
			RequestRejectedDetails: &pb.RequestRejectedDetails{
				Reason: "not enough document",
				Time:   timestamppb.New(time.Now()),
			},
		},
	}

	res := pb.CreateRequesterNotificationResponse{}

	err := its.Nhandler.CreateRequesterNotification(its.ctx, req, &res)

	its.Nil(err)
	its.T().Log(res.Msg)
}

func (its *IntTestSuite) TestCreateRequesterRequestApproved() {
	req := &pb.CreateRequesterNotificationRequest{
		From: 12345,
		To:   1234567,
		MessageOneof: &pb.CreateRequesterNotificationRequest_RequestApprovedDetails{
			RequestApprovedDetails: &pb.RequestApprovedDetails{
				Time: timestamppb.New(time.Now()),
			},
		},
	}

	res := pb.CreateRequesterNotificationResponse{}

	err := its.Nhandler.CreateRequesterNotification(its.ctx, req, &res)

	its.Nil(err)
	its.T().Log(res.Msg)
}

func (its *IntTestSuite) TestCreateRequesterValidatorQuestions() {
	req := &pb.CreateRequesterNotificationRequest{
		From: 12345,
		To:   1234567,
		MessageOneof: &pb.CreateRequesterNotificationRequest_ValidatorQuestionsDetails{
			ValidatorQuestionsDetails: &pb.ValidatorQuestionsDetails{
				ValidatorQuestionIDs: []string{"A", "B", "C", "D"},
			},
		},
	}

	res := pb.CreateRequesterNotificationResponse{}

	err := its.Nhandler.CreateRequesterNotification(its.ctx, req, &res)

	its.Nil(err)
	its.T().Log(res.Msg)
}

func (its *IntTestSuite) TestGetRequesterNotifications() {
	req := &pb.GetRequesterNotificationsRequest{
		To:     12345,
		Limit:  3,
		Offset: 2,
	}

	res := pb.GetRequesterNotificationsResponse{}

	err := its.Nhandler.GetRequesterNotifications(its.ctx, req, &res)
	its.Nil(err)

}

func (its *IntTestSuite) TestCreateValidatorNearRequest() {
	req := &pb.CreateValidatorNotificationRequest{
		From: 123456,
		To:   12345678,
		MessageOneof: &pb.CreateValidatorNotificationRequest_NearRequestFoundDetails{
			NearRequestFoundDetails: &pb.NearRequestFoundDetails{
				RequestID: "123678",
				FullName:  "asghar bibokhar",
			},
		},
	}
	res := pb.CreateValidatorNotificationResponse{}

	err := its.Nhandler.CreateValidatorNotification(its.ctx, req, &res)

	its.Nil(err)
	its.T().Log(res.Msg)
}

func (its *IntTestSuite) TestCreateValidatorElectedAsValidator() {
	req := &pb.CreateValidatorNotificationRequest{
		From: 123456,
		To:   12345678,
		MessageOneof: &pb.CreateValidatorNotificationRequest_ElectedAsValidatorDetails{
			ElectedAsValidatorDetails: &pb.ElectedAsValidatorDetails{
				RequesterFullName: "mohamad belgheisi",
			},
		},
	}
	res := pb.CreateValidatorNotificationResponse{}

	err := its.Nhandler.CreateValidatorNotification(its.ctx, req, &res)

	its.Nil(err)
	its.T().Log(res.Msg)
}

func (its *IntTestSuite) TestCreateValidatorRequsterRespond() {
	req := &pb.CreateValidatorNotificationRequest{
		From: 123456,
		To:   12345678,
		MessageOneof: &pb.CreateValidatorNotificationRequest_RequesterRespondToValidatorDetails{
			RequesterRespondToValidatorDetails: &pb.RequesterRespondToValidatorDetails{
				RequestID: "12356778",
				StoredNotifications: []*pb.StoredNotification{
					&pb.StoredNotification{
						ChatID:    11144456778,
						MessageID: 25794,
					},
					&pb.StoredNotification{
						ChatID:    234782958,
						MessageID: 692705,
					},
				},
			},
		},
	}
	res := pb.CreateValidatorNotificationResponse{}

	err := its.Nhandler.CreateValidatorNotification(its.ctx, req, &res)

	its.Nil(err)
	its.T().Log(res.Msg)
}

func (its *IntTestSuite) TestGetValidatorNotifications() {
	req := &pb.GetValidatorNotificationsRequest{
		To:     123456,
		Limit:  3,
		Offset: 1,
	}

	res := pb.GetValidatorNotificationsResponse{}

	err := its.Nhandler.GetValidatorNotifications(its.ctx, req, &res)

	its.Nil(err)
	its.T().Log(res.Notifications)
}

func (its *IntTestSuite) TearDownSuite() {
	//should clean up connection to redis
	// its.redisClient.FlushDB(its.ctx)

}

func (its *IntTestSuite) BeforeTest(suiteName, testName string) {

}

func (its *IntTestSuite) AfterTest(suiteName, testName string) {

}
