package utils

import (
	"log"
	"os"
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

type SocketWithLock struct {
	Socket *socketio.Conn
	sync.RWMutex
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func PrintDebug(format string, messagef ...interface{}) {
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		Println("debug", "", "", "", "", format, messagef...)
	}
}

func PrintInfo(message string) {
	var v []interface{}
	Println("info", "", "", "", message, "%s", v)
}

func Println(loglevel string, tenantId string, calleeServiceName string, userId string, message string, format string, messagef ...interface{}) {
	// "requestId", "guid", "tenantId", "callerServiceName", "calleeServiceName",
	// "relatedRequestId", "groupRequestId", "taskId", "userId", "entityId",
	// "objectType", "className", "method", "newTimestamp", "action",
	// "inclusiveTime", "messageCode", "instanceId", "logMessage"
	var messageTemplate string = `[` + loglevel + `] [] [] [` + tenantId + `] [Go-Notification] [` + calleeServiceName + `] [] [] [] [` + userId + `] [] [] [] [] [] [] [] [] [] [` + message + `]`
	switch loglevel {
	case "panic":
		log.Panic(messageTemplate)
		break
	case "fatal":
		log.Fatal(messageTemplate)
	case "info":
		log.Println(messageTemplate)
		break
	case "debug":
		log.Printf(`[`+loglevel+`] [] [] [`+tenantId+`] [Go-Notification] [`+calleeServiceName+`] [] [] [] [`+userId+`] [] [] [] [] [] [] [] [] [] [`+format+`]`, messagef...)
	}
}
