package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mamadeusia/NotificationSrv/domain/notification"
	"github.com/mamadeusia/NotificationSrv/entity"

	"github.com/mamadeusia/NotificationSrv/config"

	"go-micro.dev/v4/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository is our factory object
type MongoRepository struct {
	db *mongo.Database
	// message is used to store messages
	message  *mongo.Collection
	question *mongo.Collection
}

// NewMongoRepository is out factory function
func NewMongoRepository(ctx context.Context, connectionString string) (*MongoRepository, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURL()))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	notificationDatabase := client.Database("notification")
	messagesCollection := notificationDatabase.Collection("messages")
	questionsCollection := notificationDatabase.Collection("questions")
	return &MongoRepository{
		db:       notificationDatabase,
		message:  messagesCollection,
		question: questionsCollection,
	}, nil
}

func (m *MongoRepository) DropNotificationRepository(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := m.message.Drop(ctx)
	if err != nil {
		logger.Info("MONGODOMAIN: Drop database collection failed")
	}

}

// CreateMessage - is responsible for storing one instance of message entity at a time
func (m *MongoRepository) CreateNotification(ctx context.Context, notif notification.Notification) (*notification.Notification, error) {
	// TODO : this is something we have to discuss about
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	insertResult, err := m.message.InsertOne(ctx, notif)
	if err != nil {
		logger.Info("MONGODOMAIN: CreateMessage has failed with error %v", err)
		return nil, err
	}
	id, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("mongo notification : failed to convert InsertedID to primitive.ObjectID")
	}
	notif.BaseMessage.ID = id.Hex()
	return &notif, nil
}

// CreateBulkMessages - is responsible for stroring multiple instance of message entity at a time
func (m *MongoRepository) CreateBulkNotification(ctx context.Context, notifs []notification.Notification) ([]notification.Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Get().Mongo.RequestTimeOut*time.Second)
	defer cancel()

	notifInterface := make([]interface{}, len(notifs))

	for i := range notifs {
		notifInterface[i] = notifs[i]
	}
	insertResult, err := m.message.InsertMany(ctx, notifInterface)
	if err != nil {
		logger.Info("MONGODOMAIN: CreateBulkMessages has failed with error %+v", err)
		return nil, err
	}
	for i, id := range insertResult.InsertedIDs {
		objectId, ok := id.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("mongo notification : failed to convert InsertedID to primitive.ObjectID")
		}
		// TODO :: absolutely not sure about this , need to discuss about the order of the IDs
		notifs[i].BaseMessage.ID = objectId.Hex()
	}
	return notifs, nil
}

// GetMessage - is responsible for fetching data
func (m *MongoRepository) GetRequesterNotifications(ctx context.Context, to int64, limit, offset int64) ([]notification.Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(60*time.Second))
	defer cancel()
	// TODO :: if new notification create while fetching the data , it will be missed. Need to discuss about this.
	skip := offset*limit - limit
	fOpt := options.FindOptions{Limit: &limit, Skip: &skip}

	// _ = mongo.NewDeleteManyModel().Filter
	cursor, err := m.message.Find(ctx, bson.M{"basemessage.to": to}, &fOpt)
	if err != nil {
		logger.Info("MONGODOMAIN: GetMessage has failed with error %+v", err)
		return nil, err
	}
	var results []notification.Notification
	if err = cursor.All(ctx, &results); err != nil {
		logger.Info("MONGODOMAIN: GetMessage has failed with error %+v", err)
		return nil, err
	}
	return results, nil
}

// GetMessage - is responsible for fetching data
func (m *MongoRepository) GetValidatorNotifications(ctx context.Context, to int64, limit, offset int64) ([]notification.Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(60*time.Second))
	defer cancel()
	// TODO :: if new notification create while fetching the data , it will be missed. Need to discuss about this.
	skip := offset*limit - limit
	fOpt := options.FindOptions{Limit: &limit, Skip: &skip}

	cursor, err := m.message.Find(ctx, validatorNotificationFilter(to), &fOpt)
	if err != nil {
		logger.Info("MONGODOMAIN: GetMessage has failed with error %+v", err)
		return nil, err
	}
	results := []notification.Notification{}
	if err = cursor.All(ctx, &results); err != nil {
		logger.Info("MONGODOMAIN: GetMessage has failed with error %+v", err)
		return nil, err
	}
	return results, nil
}

// GetUnreadMessagesCount - is responsible for counting the documents
func (m *MongoRepository) GetUnreadRequesterNotificationCount(ctx context.Context, to int64) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(60*time.Second))
	defer cancel()
	count, err := m.message.CountDocuments(ctx, requesterCountNotificationFilter(to))
	if err != nil {
		logger.Info("MONGODOMAIN: GetUnreadMessagesCount has failed with error %v", err)
		return 0, err
	}

	return count, nil
}

func (m *MongoRepository) GetUnreadValidatorNotificationCount(ctx context.Context, to int64) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(60*time.Second))
	defer cancel()
	count, err := m.message.CountDocuments(ctx, validatorCountNotificationFilter(to))
	if err != nil {
		logger.Info("MONGODOMAIN: GetUnreadMessagesCount has failed with error %v", err)
		return 0, err
	}

	return count, nil
}

// SetMessageStatus - is responsible for set the status for documents
func (m *MongoRepository) SetNotificationStatus(ctx context.Context, id string, status entity.MessageStatus) error {
	ctx, cancel := context.WithTimeout(ctx, config.Get().Mongo.RequestTimeOut*time.Second)
	defer cancel()
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Info("MONGODOMAIN: ObjectIDFromHex has failed with error %v", err)
		return err
	}
	filter := bson.D{{Key: "_id", Value: _id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "basemessage.message_status", Value: status}}}}
	result, err := m.message.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Info("MONGODOMAIN: SetMessageStatus has failed with error %v", err)
		return err
	}

	//just for testing should be removed

	logger.Info("updated result is : %v", result)

	return nil
}
