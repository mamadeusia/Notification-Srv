package notificationEvent

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/mamadeusia/NotificationSrv/client/nats"
	"github.com/mamadeusia/NotificationSrv/config"
	"github.com/mamadeusia/NotificationSrv/entity"
	"go-micro.dev/v4/events"
	"go-micro.dev/v4/logger"

	"github.com/mamadeusia/NotificationSrv/service/notification"
	"github.com/mamadeusia/NotificationSrv/service/question"
)

const (
	groupName = "notification"
)

type Handler struct {
	prefixTopic string

	checkQuestionerValidatorCountTopics []string

	questionService     *question.QuestionService
	notificationService *notification.NotificationService

	store  events.Store
	stream events.Stream
}

func New(
	prefixTopic string,
	checkQuestionerValidatorCountTopics []string,
	store events.Store,
	stream events.Stream,
) (*Handler, error) {
	return &Handler{
		prefixTopic:                         prefixTopic,
		checkQuestionerValidatorCountTopics: checkQuestionerValidatorCountTopics,
		store:                               store,
		stream:                              stream,
	}, nil
}

// StartConsumer - nats consumer
func (n *Handler) StartConsumer(ctx context.Context) error {

	sharedPullBasedConfigurations := []nats.PullBasedTopicConfiguration{
		nats.WithPullBasedGroup(groupName),
		nats.WithPullBasedAutoAck(true),
	}
	// configure pull based configurations.
	var pullBasedTopics []*nats.PullBasedTopic

	for i, topic := range n.checkQuestionerValidatorCountTopics {

		// Important// topics be like Notification.5Second , Notification.10Minute , ...
		topicConfiguration := append(sharedPullBasedConfigurations, nats.WithPullBasedTopicString(n.prefixTopic+"."+topic+config.NatsNotificationRetryUnitString()))
		topicConsumeDelay, err := strconv.Atoi(topic)
		if err != nil {
			logger.Fatal(err)
		}
		topicConfiguration = append(topicConfiguration, nats.WithPullBasedMaxItems(10))
		topicConfiguration = append(topicConfiguration, nats.WithPullBasedDuration(time.Duration((2*topicConsumeDelay/3)*int(time.Second))))

		topicConfiguration = append(topicConfiguration, nats.WithPullBasedConsumeDelayTime(time.Duration(topicConsumeDelay*int(time.Second))))

		if i == len(n.checkQuestionerValidatorCountTopics)-1 {
			topicConfiguration = append(topicConfiguration, nats.WithPullBasedEventHandler(
				nats.EventHandler(n.checkQuestionerValidatorCountTopicsHandler_WithRejectRequestFailure())),
			)
		} else {
			topicConfiguration = append(topicConfiguration, nats.WithPullBasedEventHandler(
				nats.EventHandler(n.checkQuestionerValidatorCountTopicsHandler_WithRepulishInFailure(n.prefixTopic+"."+n.checkQuestionerValidatorCountTopics[i+1]+config.NatsNotificationRetryUnitString()))),
			)

		}

		newTopic, err := nats.NewPullBasedTopic(topicConfiguration...)
		if err != nil {
			logger.Fatal(err)
		}
		pullBasedTopics = append(pullBasedTopics, newTopic)
	}

	var natConfigurations []nats.NatsConfiguration

	natConfigurations = append(natConfigurations, nats.WithPullBasedStore(n.store))
	natConfigurations = append(natConfigurations, nats.WithPushBasedStream(n.stream))

	for _, pullBasedTopic := range pullBasedTopics {
		natConfigurations = append(natConfigurations, nats.WithPullBasedTopic(pullBasedTopic))
	}

	natsClient, err := nats.New(natConfigurations...)
	if err != nil {
		logger.Fatal(err)
	}

	return natsClient.Start(ctx)

}

// TODO error handling
func (h *Handler) checkQuestionerValidatorCountTopicsHandler_WithRepulishInFailure(nextTopic string) nats.EventHandler {
	return func(ctx context.Context, e []*events.Event) error {
		for _, event := range e {

			callbackData := entity.QuestionConfirmedReq{}
			if err := json.Unmarshal(event.Payload, &callbackData); err != nil {
				logger.Error("SERVICE::CallWithRepublish_FailureScenario unmarshalling, has failed with error : %v", err)
				continue
			}

			validatorIDs, err := h.questionService.GetQuestionerValidatorForRequestID(ctx, callbackData.RequestID, int(callbackData.MinRequiredQuestions))
			if err != nil {
				logger.Error(err)
				h.stream.Publish(nextTopic, callbackData)
				continue
			}

			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(validatorIDs), func(i, j int) { validatorIDs[i], validatorIDs[j] = validatorIDs[j], validatorIDs[i] })

			randomSelectedValidatorIDs := validatorIDs[0:callbackData.MinRequiredQuestions]

			questionIDs, err := h.questionService.GetQuestionsByValidatorIDs(ctx, callbackData.RequestID, randomSelectedValidatorIDs)
			if err != nil {
				logger.Error(err)
			}

			err = h.notificationService.CreateRequesterNotification(ctx, entity.CreateRequesterNotificationRequest{
				From: 0,
				To:   callbackData.RequesterID,
				DetailMessage: &entity.ValidatorQuestionsDetails{
					ValidatorQuestionIds: questionIDs,
				},
			})
			if err != nil {
				logger.Error(err)
			}
		}
		return nil
	}
}

func (h *Handler) checkQuestionerValidatorCountTopicsHandler_WithRejectRequestFailure() nats.EventHandler {
	return func(ctx context.Context, e []*events.Event) error {
		for _, event := range e {

			callbackData := entity.QuestionConfirmedReq{}
			if err := json.Unmarshal(event.Payload, &callbackData); err != nil {
				logger.Error("SERVICE::CallWithRepublish_FailureScenario unmarshalling, has failed with error : %v", err)
				continue
			}

			validatorIDs, err := h.questionService.GetQuestionerValidatorForRequestID(ctx, callbackData.RequestID, int(callbackData.MinRequiredQuestions))
			if err != nil {
				logger.Error(err)
				err = h.notificationService.CreateRequesterNotification(ctx, entity.CreateRequesterNotificationRequest{
					From: 0,
					To:   callbackData.RequesterID,
					DetailMessage: &entity.RequestRejectedDetails{
						Reason: "Because We didn't have enough validators",
						Time:   time.Now(),
					},
				})
				if err != nil {
					logger.Error(err)
				}

			}

			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(validatorIDs), func(i, j int) { validatorIDs[i], validatorIDs[j] = validatorIDs[j], validatorIDs[i] })

			randomSelectedValidatorIDs := validatorIDs[0:callbackData.MinRequiredQuestions]

			questionIDs, err := h.questionService.GetQuestionsByValidatorIDs(ctx, callbackData.RequestID, randomSelectedValidatorIDs)
			if err != nil {
				logger.Error(err)
			}

			err = h.notificationService.CreateRequesterNotification(ctx, entity.CreateRequesterNotificationRequest{
				From: 0,
				To:   callbackData.RequesterID,
				DetailMessage: &entity.ValidatorQuestionsDetails{
					ValidatorQuestionIds: questionIDs,
				},
			})
			if err != nil {
				logger.Error(err)
			}
		}
		return nil
	}
}
