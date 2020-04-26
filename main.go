package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ivegotwings/mdm-ui-go/redis"

	socketio "github.com/googollee/go-socket.io"
)

type Config struct {
	Redis struct {
		Host string
		Port string
	}
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	_ = json.Unmarshal([]byte(byteValue), &config)
	return config
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	var config Config = LoadConfiguration("config.json")
	b, err := json.Marshal(config)
	fmt.Println("Redis Config-")
	fmt.Println(string(b))

	opts := make(map[string]string)
	opts["host"] = config.Redis.Host
	opts["port"] = config.Redis.Port
	server.SetAdaptor(redis.Redis(opts))

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			log.Println("emit:", so.Emit("chat message", msg))
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
