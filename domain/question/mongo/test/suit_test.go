package test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/mamadeusia/NotificationSrv/config"
	"github.com/mamadeusia/NotificationSrv/domain/question"
	questionMongo "github.com/mamadeusia/NotificationSrv/domain/question/mongo"
	"github.com/mamadeusia/NotificationSrv/entity"
	"github.com/stretchr/testify/suite"
)

// integration test suite
type IntTestSuite struct {
	suite.Suite
	questionRepo        *questionMongo.MongoRepository
	ctx                 context.Context
	validator1toreq1_id []string
	validator2toreq2_id []string
	validator3toreq2_id []string
	request1Questions   []question.Question
	request2Questions   []question.Question
}

var (
	validator1Id = int64(12345)
	validator2Id = int64(23456)
	validator3Id = int64(56789)
	requester1Id = int64(34567)
	requester2Id = int64(45678)
	request_id1  = uuid.New().String()
	request_id2  = uuid.New().String()
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
	questionRepo, err := questionMongo.NewMongoRepository(its.ctx, config.MongoURL())
	if err != nil {
		its.FailNow("MAINROUTINE::NewMongoRepository:: has failed with error", err)
	}
	its.questionRepo = questionRepo
	its.questionRepo.DropQuestionRepository(its.ctx)

}

func (its *IntTestSuite) TestCreateBulkQuestions() {
	its.request1Questions = []question.Question{
		question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    request_id1,
				ValidatorID:  validator1Id,
				RequesterID:  requester1Id,
				QuestionType: entity.DescriptiveQuestion,
			},
			QuestionDetails: &entity.DescriptiveQuestionDetails{
				Question: "salam chetori",
			},
		},
		question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    request_id1,
				ValidatorID:  validator1Id,
				RequesterID:  requester1Id,
				QuestionType: entity.DescriptiveQuestion,
			},
			QuestionDetails: &entity.DescriptiveQuestionDetails{
				Question: "bidari?",
			},
		},
		question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    request_id1,
				ValidatorID:  validator1Id,
				RequesterID:  requester1Id,
				QuestionType: entity.DescriptiveQuestion,
			},
			QuestionDetails: &entity.DescriptiveQuestionDetails{
				Question: "che khabar",
			},
		},
		question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    request_id1,
				ValidatorID:  validator1Id,
				RequesterID:  requester1Id,
				QuestionType: entity.MultipleChoiceQuestion,
			},
			QuestionDetails: &entity.MultipleChoiceQuestionDetails{
				Question: "chi mizani ",
				Choises:  []string{"gol", "sigar", "abe karafs"},
			},
		},
	}

	its.request2Questions = []question.Question{
		question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    request_id2,
				ValidatorID:  validator2Id,
				RequesterID:  requester2Id,
				QuestionType: entity.MultipleChoiceQuestion,
			},
			QuestionDetails: &entity.MultipleChoiceQuestionDetails{
				Question: "chi dostdari ",
				Choises:  []string{"hichi", "abe havj", "moz"},
			},
		},
		question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    request_id2,
				ValidatorID:  validator2Id,
				RequesterID:  requester2Id,
				QuestionType: entity.DescriptiveQuestion,
			},
			QuestionDetails: &entity.DescriptiveQuestionDetails{
				Question: "dige che khabar",
			},
		},
		question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    request_id2,
				ValidatorID:  validator2Id,
				RequesterID:  requester2Id,
				QuestionType: entity.DescriptiveQuestion,
			},
			QuestionDetails: &entity.DescriptiveQuestionDetails{
				Question: "khob migofti",
			},
		},
		question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    request_id2,
				ValidatorID:  validator3Id,
				RequesterID:  requester2Id,
				QuestionType: entity.MultipleChoiceQuestion,
			},
			QuestionDetails: &entity.MultipleChoiceQuestionDetails{
				Question: "che timi ",
				Choises:  []string{"perspolis", "esteghlal", "sepahan"},
			},
		},
		question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    request_id2,
				ValidatorID:  validator3Id,
				RequesterID:  requester2Id,
				QuestionType: entity.DescriptiveQuestion,
			},
			QuestionDetails: &entity.DescriptiveQuestionDetails{
				Question: "chikar mikoni",
			},
		},
		question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    request_id2,
				ValidatorID:  validator3Id,
				RequesterID:  requester2Id,
				QuestionType: entity.DescriptiveQuestion,
			},
			QuestionDetails: &entity.DescriptiveQuestionDetails{
				Question: "koja miri",
			},
		},
	}

	err := its.questionRepo.CreateBulkQuestions(its.ctx, its.request1Questions)
	its.Nil(err)

	err = its.questionRepo.CreateBulkQuestions(its.ctx, its.request2Questions)
	its.Nil(err)

	for i := range its.request1Questions {
		q := its.request1Questions[i]
		its.validator1toreq1_id = append(its.validator1toreq1_id, q.BaseQuestion.ID)
	}

	for i := range its.request2Questions {
		q := its.request2Questions[i]
		if q.BaseQuestion.ValidatorID == validator2Id {
			its.validator2toreq2_id = append(its.validator2toreq2_id, q.BaseQuestion.ID)
		} else {
			its.validator3toreq2_id = append(its.validator3toreq2_id, q.BaseQuestion.ID)
		}
	}

}

func check_id_equality(question_id []string, expected_id []string) bool {
	return reflect.DeepEqual(question_id, expected_id)
}

func check_id(question_id []string, expected_id []string, limit int64, offset int64) bool {
	if offset*limit > int64(len(expected_id)) {
		return check_id_equality(question_id, expected_id[offset*limit-limit:])
	} else {
		return check_id_equality(question_id, expected_id[offset*limit-limit:offset*limit])
	}
}

func (its *IntTestSuite) TestGetValidatorQuestionByRequest() {
	questions_id, err := its.questionRepo.GetValidatorQuestionsByRequestID(its.ctx, validator1Id, request_id1, 2, 1)
	its.Nil(err)
	its.True(check_id(questions_id, its.validator1toreq1_id, 2, 1))

	questions_id, err = its.questionRepo.GetValidatorQuestionsByRequestID(its.ctx, validator1Id, request_id1, 2, 2)
	its.Nil(err)
	its.True(check_id(questions_id, its.validator1toreq1_id, 2, 2))

	questions_id, err = its.questionRepo.GetValidatorQuestionsByRequestID(its.ctx, validator1Id, request_id1, 5, 1)
	its.Nil(err)
	its.True(check_id(questions_id, its.validator1toreq1_id, 5, 1))

	questions_id, err = its.questionRepo.GetValidatorQuestionsByRequestID(its.ctx, validator2Id, request_id2, 2, 1)
	its.Nil(err)
	its.True(check_id(questions_id, its.validator2toreq2_id, 2, 1))

	questions_id, err = its.questionRepo.GetValidatorQuestionsByRequestID(its.ctx, validator2Id, request_id2, 2, 2)
	its.Nil(err)
	its.True(check_id(questions_id, its.validator2toreq2_id, 2, 2))

	questions_id, err = its.questionRepo.GetValidatorQuestionsByRequestID(its.ctx, validator3Id, request_id2, 2, 1)
	its.Nil(err)
	its.True(check_id(questions_id, its.validator3toreq2_id, 2, 1))

	questions_id, err = its.questionRepo.GetValidatorQuestionsByRequestID(its.ctx, validator3Id, request_id2, 2, 2)
	its.Nil(err)
	its.True(check_id(questions_id, its.validator3toreq2_id, 2, 2))
}

func (its *IntTestSuite) TestGetQuestionerValidatorsIDsByRequestID() {
	ids, err := its.questionRepo.GetQuestionerValidatorsIDsByRequestID(its.ctx, request_id1)
	its.Nil(err)
	its.Equal(ids, []int64{validator1Id})

	ids, err = its.questionRepo.GetQuestionerValidatorsIDsByRequestID(its.ctx, request_id2)
	its.Nil(err)
	its.Equal(ids, []int64{validator2Id, validator3Id})
}

func (its *IntTestSuite) TestGetQuestionerValidatorCountByRequestID() {
	cnt, err := its.questionRepo.GetQuestionerValidatorCountByRequestID(its.ctx, request_id1)
	its.Nil(err)
	its.Equal(cnt, int64(len(its.request1Questions)))

	cnt, err = its.questionRepo.GetQuestionerValidatorCountByRequestID(its.ctx, request_id2)
	its.Nil(err)
	its.Equal(cnt, int64(len(its.request2Questions)))
}

func (its *IntTestSuite) TestGetQuestionByID() {
	// TODO : check return result
	for _, q := range its.request1Questions {
		question, err := its.questionRepo.GetQuestionByID(its.ctx, q.BaseQuestion.ID)
		its.Nil(err)
		its.True(reflect.DeepEqual(*question, q))
	}

	for _, q := range its.request2Questions {
		question, err := its.questionRepo.GetQuestionByID(its.ctx, q.BaseQuestion.ID)
		its.Nil(err)
		its.True(reflect.DeepEqual(*question, q))
	}

}

func (its *IntTestSuite) TestUpdateQuestion() {
	answer := question.Question{
		BaseQuestion: &entity.BaseQuestion{
			RequestID:    request_id1,
			ValidatorID:  validator1Id,
			RequesterID:  requester1Id,
			QuestionType: entity.DescriptiveQuestion,
		},
		QuestionDetails: &entity.DescriptiveQuestionDetails{
			Question: "bidari?",
			Answer:   "na khabam",
		},
	}
	err := its.questionRepo.UpdateQuestion(its.ctx, its.request1Questions[1].BaseQuestion.ID, answer)
	its.Nil(err)

	q, err := its.questionRepo.GetQuestionByID(its.ctx, its.request1Questions[1].BaseQuestion.ID)
	its.Nil(err)
	its.True(reflect.DeepEqual(q.QuestionDetails, answer.QuestionDetails))

	answer = question.Question{
		BaseQuestion: &entity.BaseQuestion{
			RequestID:    request_id2,
			ValidatorID:  validator2Id,
			RequesterID:  requester2Id,
			QuestionType: entity.MultipleChoiceQuestion,
		},
		QuestionDetails: &entity.MultipleChoiceQuestionDetails{
			Question:    "chi dostdari ",
			Choises:     []string{"hichi", "abe havj", "moz"},
			AnswerIndex: 1,
		},
	}
	err = its.questionRepo.UpdateQuestion(its.ctx, its.request2Questions[0].BaseQuestion.ID, answer)
	its.Nil(err)

	q, err = its.questionRepo.GetQuestionByID(its.ctx, its.request2Questions[0].BaseQuestion.ID)
	its.Nil(err)
	its.True(reflect.DeepEqual(q.QuestionDetails, answer.QuestionDetails))
}
