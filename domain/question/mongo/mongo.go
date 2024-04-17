package mongo

import (
	"context"
	"log"
	"time"

	"github.com/mamadeusia/NotificationSrv/config"
	"github.com/mamadeusia/NotificationSrv/domain/question"
	"github.com/mamadeusia/NotificationSrv/entity"
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

func (m *MongoRepository) DropQuestionRepository(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := m.question.Drop(ctx)
	if err != nil {
		logger.Info("MONGODOMAIN: Drop database collection failed")
	}

}

func removeDuplicates(slice []int64) []int64 {
	// Create a map to store unique elements
	seen := make(map[int64]bool)
	result := []int64{}

	// Loop through the slice, adding elements to the map if they haven't been seen before
	for _, val := range slice {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}

// CreateBulkQuestions - is responsible for storing multiple instance of question entity
func (m *MongoRepository) CreateBulkQuestions(ctx context.Context, questions []question.Question) error {
	ctx, cancel := context.WithTimeout(ctx, config.Get().Mongo.RequestTimeOut*time.Second)
	defer cancel()
	questionInterface := make([]interface{}, len(questions))

	for i := range questions {
		questionInterface[i] = questions[i]
	}
	result, err := m.question.InsertMany(ctx, questionInterface)
	if err != nil {
		logger.Info("MONGODOMAIN: CreateBulkQuestions has failed with error %+v", err)
		return err
	}

	for index := range result.InsertedIDs {
		id := result.InsertedIDs[index].(primitive.ObjectID)
		questions[index].BaseQuestion.ID = id.Hex()
	}

	return nil
}

func (m *MongoRepository) UpdateQuestion(ctx context.Context, questionID string, question question.Question) error {
	ctx, cancel := context.WithTimeout(ctx, config.Get().Mongo.RequestTimeOut*time.Second)
	defer cancel()
	if question.BaseQuestion.QuestionType == entity.DescriptiveQuestion {
		qDestails, err := question.QuestionDetails.GetQuestionDetails()
		if err != nil {
			return err
		}

		answer, err := GetValueMap[string]("answer", qDestails)
		if err != nil {
			return err
		}
		_id, err := primitive.ObjectIDFromHex(questionID)
		if err != nil {
			logger.Info("MONGODOMAIN: ObjectIDFromHex has failed with error %v", err)
			return err
		}
		filter := bson.D{{Key: "_id", Value: _id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "questiondetails.answer", Value: answer}}}}
		if _, err := m.question.UpdateOne(ctx, filter, update); err != nil {
			logger.Info("MONGODOMAIN: UpdateQuestion has failed with error %+v", err)
			return err
		}

	} else if question.BaseQuestion.QuestionType == entity.MultipleChoiceQuestion {
		qDestails, err := question.QuestionDetails.GetQuestionDetails()
		if err != nil {
			return err
		}

		answer, err := GetValueMap[int]("answer_index", qDestails)
		if err != nil {
			return err
		}
		_id, err := primitive.ObjectIDFromHex(questionID)
		if err != nil {
			logger.Info("MONGODOMAIN: ObjectIDFromHex has failed with error %v", err)
			return err
		}
		filter := bson.D{{Key: "_id", Value: _id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "questiondetails.answer_index", Value: answer}}}}
		if _, err := m.question.UpdateOne(context.TODO(), filter, update); err != nil {
			logger.Info("MONGODOMAIN: UpdateQuestion has failed with error %+v", err)
			return err
		}
	}

	return nil
}

// GetQuestions - is responsible for fetching data
func (m *MongoRepository) GetValidatorQuestionsByRequestID(ctx context.Context, validatorID int64, requestID string, limit, offset int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Get().Mongo.RequestTimeOut*time.Second)
	defer cancel()

	skip := offset*limit - limit
	fOpt := options.FindOptions{Limit: &limit, Skip: &skip}
	// opts := options.Find().SetProjection(bson.D{{Key: "id", Value: 1}})

	cursor, err := m.question.Find(ctx, validatorQuestionFilter(requestID, validatorID), &fOpt)
	if err != nil {
		logger.Info("MONGODOMAIN: GetQuestions has failed with error %v", err)
		return nil, err
	}

	var questions []question.Question
	if err := cursor.All(ctx, &questions); err != nil {
		logger.Info("MONGODOMAIN: Decoding result to question entity has failed with error %v", err)
		return nil, err
	}
	var ids []string
	for _, q := range questions {
		ids = append(ids, q.BaseQuestion.ID)
	}
	return ids, nil
}

// GetValidatorsFromQuestionsByRequesterID - is responsible for extracting validators from questions
func (m *MongoRepository) GetQuestionerValidatorsIDsByRequestID(ctx context.Context, requestID string) ([]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Get().Mongo.RequestTimeOut*time.Second)
	defer cancel()
	// findOptions := options.Find()
	// Sort by `created_at` field descending
	// findOptions.SetSort(bson.D{{"created_at", -1}})

	cursor, err := m.question.Find(ctx, bson.M{"basequestion.request_id": requestID})
	if err != nil {
		logger.Info("MONGODOMAIN: GetQuestions has failed with error %v", err)
		return nil, err
	}

	var questions []question.Question
	if err := cursor.All(ctx, &questions); err != nil {
		logger.Info("MONGODOMAIN: Decoding result to question entity has failed with error %v", err)
		return nil, err
	}

	var validators []int64

	for _, q := range questions {
		validators = append(validators, q.BaseQuestion.ValidatorID)
	}

	return removeDuplicates(validators), nil
}

// GetQuestionsCount -
func (m *MongoRepository) GetQuestionerValidatorCountByRequestID(ctx context.Context, requestID string) (int64, error) {
	itemCount, err := m.question.CountDocuments(ctx, bson.M{"basequestion.request_id": requestID})
	if err != nil {
		logger.Info("MONGODOMAIN: GetQuestionsCount has failed with error %v", err)
		return 0, err
	}

	return itemCount, nil
}

func (m *MongoRepository) GetQuestionByID(ctx context.Context, id string) (*question.Question, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Get().Mongo.RequestTimeOut*time.Second)
	defer cancel()
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Info("MONGODOMAIN: ObjectIDFromHex has failed with error %v", err)
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: _id}}
	cursor := m.question.FindOne(ctx, filter)

	var question question.Question
	if err := cursor.Decode(&question); err != nil {
		logger.Info("MONGODOMAIN: Decoding result to question entity has failed with error %v", err)
		return nil, err
	}
	return &question, nil
}
