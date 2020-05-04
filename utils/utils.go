package utils

import (
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

type SocketWithLock struct {
	Socket *socketio.Conn
	Lock   *sync.Mutex
}
