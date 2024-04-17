package main

import (
	"context"
	"fmt"

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

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "notificationsrv"
	version = "latest"
)

func main() {

	// Load conigurations
	if err := config.Load(); err != nil {
		logger.Fatal("MAINROUTINE::LoadConfig:: has failed with error : %+v", err)
	}

	fmt.Println(config.Get().Mongo)
	ctx := context.Background()

	notifRepo, err := notifMongo.NewMongoRepository(ctx, config.MongoURL())
	if err != nil {
		logger.Fatal("MAINROUTINE::NewMongoRepository:: has failed with error : %+v", err)
	}

	questionRepo, err := questionMongo.NewMongoRepository(ctx, config.MongoURL())
	if err != nil {
		logger.Fatal("MAINROUTINE::NewMongoRepository:: has failed with error : %+v", err)
	}
	stream, err := natsjs.NewStream(
		natsjs.Address(config.NatsURL()),
		// natsjs.NkeyConfig(config.NatsNkey()),
	)
	if err != nil {
		logger.Fatal(err)
	}

	store, err := natsjs.NewStore(
		natsjs.Address(config.NatsURL()),
		// natsjs.NkeyConfig(config.NatsNkey()),
	)

	if err != nil {
		logger.Error(err)
	}

	notificationService := notifService.NewNotificationService(notifRepo, stream)
	quesitonService := questionService.NewQuestionService(questionRepo, stream)
	// Create service
	srv := micro.NewService(
		// micro.Server(grpcs.NewServer()),
		// micro.Client(grpcc.NewClient()),
		micro.Context(ctx),
	)
	srv.Init(
		micro.Name(service),
		micro.Version(version),
	)

	Nhandler := notificationGrpc.New(notificationService)
	Qhandler := questionGrpc.New(quesitonService)
	EventHandler, err := notificationEvent.New(config.NatsNotificationPrefix(), config.NatsNotificationRetryDurations(), store, stream)
	if err != nil {
		logger.Fatal(err)
	}

	// Register handler
	if err := pb.RegisterNotificationSrvHandler(srv.Server(), Nhandler); err != nil {
		logger.Fatal(err)
	}
	if err := pb.RegisterQuestionSrvHandler(srv.Server(), Qhandler); err != nil {
		logger.Fatal(err)
	}

	if err := EventHandler.StartConsumer(ctx); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

}
