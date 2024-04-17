package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

func validatorQuestionFilter(requestID string, validatorId int64) interface{} {
	return bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "basequestion.request_id", Value: requestID}},
				bson.D{{Key: "basequestion.validator_id", Value: validatorId}},
			},
		},
	}
}
