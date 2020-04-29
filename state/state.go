package state

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

var conn redis.Conn

func Connect(opts map[string]string) error {
	fmt.Println("state- redis config", opts)
	var ok bool
	var host, port string
	host, ok = opts["host"]
	if !ok {
		host = "127.0.0.1"
	}
	port, ok = opts["port"]
	if !ok {
		port = "6379"
	}
	conn, err := redis.Dial("tcp", host+":"+port)

	if err != nil {
		return errors.New("state- unable to connect to redis")
	}
	fmt.Println(conn)
	return nil
}
