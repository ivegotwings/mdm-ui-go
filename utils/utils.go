package utils

import (
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

type SocketWithLock struct {
	Socket *socketio.Conn
	Lock   *sync.Mutex
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func Log(tenantId string, calleeServiceName string, userId string, message string) {
	// "requestId", "guid", "tenantId", "callerServiceName", "calleeServiceName",
	// "relatedRequestId", "groupRequestId", "taskId", "userId", "entityId",
	// "objectType", "className", "method", "newTimestamp", "action",
	// "inclusiveTime", "messageCode", "instanceId", "logMessage"
	//var messageTemplate string = `[] [] [` + tenantId + `] [Go-Notification] [` + calleeServiceName + `] [] [] [] [` + userId + `] [] [] [] [] [] [] [] [] [] [` + message + `]`
}
