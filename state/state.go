package state

import (
	"errors"

	"github.com/garyburd/redigo/redis"
	"github.com/ivegotwings/mdm-ui-go/utils"
)

var conn redis.Conn

func Connect(opts map[string]string) error {
	var ok bool
	var err error
	var host, port string
	host, ok = opts["host"]
	if !ok {
		host = "127.0.0.1"
	}
	port, ok = opts["port"]
	if !ok {
		port = "6379"
	}
	conn, err = redis.Dial("tcp", host+":"+port)

	if err != nil {
		utils.PrintInfo("state- unable to connect to redis")
		return errors.New("state- unable to connect to redis")
	}
	return nil
}

func Conn() *redis.Conn {
	return &conn
}
