package notification

import (
	"errors"

	"github.com/mamadeusia/NotificationSrv/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	BaseMessage    *entity.BaseMessage
	MessageDetails entity.DetailMessage
}

func (n *Notification) UnmarshalBSON(b []byte) error {
	notficationTemp := new(struct {
		ID primitive.ObjectID  `json:"id" bson:"_id"`
		Bm *entity.BaseMessage `bson:"basemessage"`
		Md bson.Raw            `bson:"messagedetails"`
	})
	if err := bson.Unmarshal(b, notficationTemp); err != nil {
		return err
	}
	n.BaseMessage = notficationTemp.Bm
	n.BaseMessage.ID = notficationTemp.ID.Hex()
	switch notficationTemp.Bm.MessageType {
	case entity.AdminAskQuestion:
		a := &entity.AdminAskQuestionDetails{}
		err := bson.Unmarshal(notficationTemp.Md, a)
		if err != nil {
			return err
		}
		n.MessageDetails = a
	case entity.RequestRejected:
		r := &entity.RequestRejectedDetails{}
		err := bson.Unmarshal(notficationTemp.Md, r)
		if err != nil {
			return err
		}
		n.MessageDetails = r
	case entity.RequestApproved:
		r := &entity.RequestApprovedDetails{}
		err := bson.Unmarshal(notficationTemp.Md, r)
		if err != nil {
			return err
		}
		n.MessageDetails = r
	case entity.ValidatorQuestions:
		v := &entity.ValidatorQuestionsDetails{}
		err := bson.Unmarshal(notficationTemp.Md, v)
		if err != nil {
			return err
		}
		n.MessageDetails = v
	case entity.NearRequestFound:
		a := &entity.NearRequestFoundDetails{}
		err := bson.Unmarshal(notficationTemp.Md, a)
		if err != nil {
			return err
		}
		n.MessageDetails = a
	case entity.ElectedAsValidator:
		a := &entity.ElectedAsValidatorDetails{}
		err := bson.Unmarshal(notficationTemp.Md, a)
		if err != nil {
			return err
		}
		n.MessageDetails = a
	case entity.RequesterRespondToValidator:
		a := &entity.RequesterRespondToValidatorDetails{}
		err := bson.Unmarshal(notficationTemp.Md, a)
		if err != nil {
			return err
		}
		n.MessageDetails = a
	default:
		return errors.New("invalid notification type")

	}

	return nil
}
