package question

import (
	"errors"

	"github.com/mamadeusia/NotificationSrv/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	BaseQuestion    *entity.BaseQuestion
	QuestionDetails entity.DetailQuestion
}

func (q *Question) UnmarshalBSON(b []byte) error {
	notficationTemp := new(struct {
		ID primitive.ObjectID   `json:"id" bson:"_id"`
		Bm *entity.BaseQuestion `bson:"basequestion"`
		Md bson.Raw             `bson:"questiondetails"`
	})
	if err := bson.Unmarshal(b, notficationTemp); err != nil {
		return err
	}
	q.BaseQuestion = notficationTemp.Bm
	q.BaseQuestion.ID = notficationTemp.ID.Hex()
	switch notficationTemp.Bm.QuestionType {
	case entity.DescriptiveQuestion:
		a := &entity.DescriptiveQuestionDetails{}
		err := bson.Unmarshal(notficationTemp.Md, a)
		if err != nil {
			return err
		}
		q.QuestionDetails = a
	case entity.MultipleChoiceQuestion:
		r := &entity.MultipleChoiceQuestionDetails{}
		err := bson.Unmarshal(notficationTemp.Md, r)
		if err != nil {
			return err
		}
		q.QuestionDetails = r
	default:
		return errors.New("invalid question type")
	}

	return nil
}
