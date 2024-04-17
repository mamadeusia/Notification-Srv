package mongo

import (
	"github.com/mamadeusia/NotificationSrv/entity"
	"go.mongodb.org/mongo-driver/bson"
)

func requesterNotificationFilter(to int64) interface{} {
	return bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "basemessage.to", Value: to}},
				bson.D{{Key: "$or", Value: bson.A{
					bson.D{{Key: "basemessage.message_type", Value: string(entity.AdminAskQuestion)}},
					bson.D{{Key: "basemessage.message_type", Value: string(entity.RequestRejected)}},
					bson.D{{Key: "basemessage.message_type", Value: string(entity.RequestApproved)}},
					bson.D{{Key: "basemessage.message_type", Value: string(entity.ValidatorQuestions)}},
				}}},
			},
		},
	}
}

func requesterCountNotificationFilter(to int64) interface{} {
	return bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "basemessage.to", Value: to}},
				bson.D{{Key: "$or", Value: bson.A{
					bson.D{{Key: "basemessage.message_type", Value: string(entity.AdminAskQuestion)}},
					bson.D{{Key: "basemessage.message_type", Value: string(entity.RequestRejected)}},
					bson.D{{Key: "basemessage.message_type", Value: string(entity.RequestApproved)}},
					bson.D{{Key: "basemessage.message_type", Value: string(entity.ValidatorQuestions)}},
				}}},
				bson.D{{Key: "basemessage.message_status", Value: entity.UnRead}},
			},
		},
	}
}

func validatorNotificationFilter(to int64) interface{} {
	return bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "basemessage.to", Value: to}},
				bson.D{{Key: "$or", Value: bson.A{
					bson.D{{Key: "basemessage.message_type", Value: string(entity.NearRequestFound)}},
					bson.D{{Key: "basemessage.message_type", Value: string(entity.ElectedAsValidator)}},
					bson.D{{Key: "basemessage.message_type", Value: string(entity.RequesterRespondToValidator)}},
				}}},
			},
		},
	}
}

func validatorCountNotificationFilter(to int64) interface{} {
	return bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "basemessage.to", Value: to}},
				bson.D{{Key: "$or", Value: bson.A{
					bson.D{{Key: "basemessage.message_type", Value: string(entity.NearRequestFound)}},
					bson.D{{Key: "basemessage.message_type", Value: string(entity.ElectedAsValidator)}},
					bson.D{{Key: "basemessage.message_type", Value: string(entity.RequesterRespondToValidator)}},
				}}},
				bson.D{{Key: "basemessage.message_status", Value: entity.UnRead}},
			},
		},
	}
}
