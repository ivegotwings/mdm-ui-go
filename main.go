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
func notifyRouterWrapper(redisBroadCastAdaptor *redis.Broadcast) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			var message []interface{}
			err = json.Unmarshal(body, &message)
			if err != nil {
				fmt.Println("ERR", err)
				log.Println("notify error in processing body", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Println(message)
			redisBroadCastAdaptor.Send(nil, "testroom", "event:notification", message...)
			fmt.Fprint(w, "POST done")
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
		w.Header().Set("Server", "A Go Web Server")
		w.WriteHeader(200)
	}

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
	//notifiy channel
	redisBroadCastAdaptor := redis.Redis(opts)

	server.OnConnect("", func(so socketio.Conn) error {
		so.SetContext("")
		err := redisBroadCastAdaptor.Join("testroom", so)
		if err != nil {
			fmt.Println("Redis BroadCastManager- Failure to connect", err)
		}
		fmt.Println("connected:", so.ID())
		log.Println("connected:", so.ID())

		return nil
	})
	server.OnError("error", func(so socketio.Conn, err error) {
		log.Println("error:", err)
	})
	server.OnEvent("/", "event:adduser", func(so socketio.Conn, msg string) {
		fmt.Println("event:adduser", msg)
		var _userInfo interface{}
		err := json.Unmarshal([]byte(msg), &_userInfo)
		if err != nil {
			fmt.Println("error processing event:adduser")
		} else {
			userInfo, ok := _userInfo.(map[string]interface{})
			fmt.Println("debug ", userInfo, ok)
			if ok {
				//join user room
				user_room := "socket_room_tenant_" + userInfo["tenantId"].(string) + "_user_" + userInfo["userId"].(string)
				err = redisBroadCastAdaptor.Join(user_room, so)
				//join tenant room
				tenant_room := "socket_room_tenant_" + userInfo["tenantId"].(string)
				err = redisBroadCastAdaptor.Join(tenant_room, so)
				if err != nil {
					fmt.Println("Redis BroadCastManager- Failure to connect", err)
				} else {
					fmt.Println("adding new user to rooms", user_room, tenant_room)
				}
			}
		}

	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.HandlerFunc(baseRouter))
	http.Handle("/notify", http.HandlerFunc(notifyRouterWrapper(redisBroadCastAdaptor)))

	fmt.Println("Serving at localhost:5007...")
	log.Println("Serving at localhost:5007...")
	log.Fatal(http.ListenAndServe(":5007", nil))
}
