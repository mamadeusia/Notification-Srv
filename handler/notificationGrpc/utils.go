package notificationGrpc

import (
	"errors"
	"time"

	"github.com/mamadeusia/NotificationSrv/entity"
)

func GetValueMap[V ~[]string | string | int32 | int64 | time.Time | []entity.StoredMessage](key string, m map[string]interface{}) (V, error) {
	var noop V
	val, ok := m[key]
	if !ok {
		return noop, errors.New("key not found in the message details")
	}
	output, ok := val.(V)
	if !ok {
		return noop, errors.New("value type not supported")
	}
	return output, nil
}
