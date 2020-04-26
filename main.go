package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
	"github.com/ivegotwings/mdm-ui-go/redis"
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

func baseRouter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "A Go Web Server")
	w.WriteHeader(200)
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	var config Config = LoadConfiguration("config.json")
	b, err := json.Marshal(config)
	fmt.Println("Redis Config-")
	fmt.Println(string(b))

	opts := make(map[string]string)
	opts["host"] = config.Redis.Host
	opts["port"] = config.Redis.Port

	redisBroadCastAdaptor := redis.Redis(opts)

	server.OnConnect("", func(so socketio.Conn) error {
		so.SetContext("")
		err := redisBroadCastAdaptor.Join("testroom", so)
		if err != nil {
			fmt.Println("Redis BroadCastManager- Failure to connect", err)
		}
		fmt.Println("connected:", so.ID())
		log.Println("connected:", so.ID())

		//		so.Join("chat")
		return nil
	})
	server.OnError("error", func(so socketio.Conn, err error) {
		log.Println("error:", err)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.HandlerFunc(baseRouter))
	fmt.Println("Serving at localhost:5007...")
	log.Println("Serving at localhost:5007...")
	log.Fatal(http.ListenAndServe(":5007", nil))
}
