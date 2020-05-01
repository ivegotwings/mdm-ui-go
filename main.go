package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/ivegotwings/mdm-ui-go/connection"
	"github.com/ivegotwings/mdm-ui-go/moduleversion"
	"github.com/ivegotwings/mdm-ui-go/notification"
	"github.com/ivegotwings/mdm-ui-go/state"
)

type Config struct {
	Redis struct {
		Host string
		Port string
	}
	NotificationInterval uint
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	if err != nil {
		log.Println(err.Error())
	}
	_ = json.Unmarshal([]byte(byteValue), &config)
	return config
}

func baseRouter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "A Go Web Server")
	w.WriteHeader(200)
}

var redisBroadCastAdaptor connection.Broadcast

func main() {
	runtime.GOMAXPROCS(100)
	fmt.Println(runtime.GOMAXPROCS(20))
	log.SetOutput(ioutil.Discard)
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	var config Config = LoadConfiguration("config.json")
	b, err := json.Marshal(config)
	log.Println("Redis Config-")
	log.Println(string(b))
	//pre load the map once
	moduleversion.LoadDomainMap()

	opts := make(map[string]string)
	opts["host"] = config.Redis.Host
	opts["port"] = config.Redis.Port
	//notifiy channel
	redisBroadCastAdaptor = *connection.Redis(opts)
	//state channel
	err = state.Connect(opts)
	if err != nil {
		panic(err)
	}
	var ticker *time.Ticker = time.NewTicker(time.Duration(config.NotificationInterval) * time.Millisecond)
	var quit = make(chan struct{})
	notification.SetRedisBroadCastAdaptor(&redisBroadCastAdaptor)
	go notification.NotificationScheduler(ticker, quit)

	server.OnConnect("", func(so socketio.Conn) error {
		so.SetContext("")
		err := redisBroadCastAdaptor.Join("testroom", so)
		if err != nil {
			log.Println("Redis BroadCastManager- Failure to connect", err)
		}
		log.Println("connected:", so.ID())
		log.Println("connected:", so.ID())

		return nil
	})
	server.OnError("error", func(so socketio.Conn, err error) {
		log.Println("error:", err)
	})

	server.OnEvent("/", "event:adduser", func(so socketio.Conn, msg string) {
		log.Println("event:adduser", msg)
		var _userInfo interface{}
		err := json.Unmarshal([]byte(msg), &_userInfo)
		if err != nil {
			log.Println("error processing event:adduser")
		} else {
			userInfo, ok := _userInfo.(map[string]interface{})
			log.Println("debug ", userInfo, ok)
			if ok {
				//join user room
				user_room := "socket_conn_room_tenant_" + userInfo["tenantId"].(string) + "_user_" + userInfo["userId"].(string)
				err = redisBroadCastAdaptor.Join(user_room, so)
				//join tenant room
				tenant_room := "socket_conn_room_tenant_" + userInfo["tenantId"].(string)
				err = redisBroadCastAdaptor.Join(tenant_room, so)

				log.Println(user_room, tenant_room)
				if err != nil {
					log.Println("Redis BroadCastManager- Failure to connect", err)
				} else {
					log.Println("adding new user to rooms", user_room, tenant_room)
					so.Emit("event:message", _userInfo)
				}
			}
		}

	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.HandlerFunc(baseRouter))
	notificationHandler := notification.NotificationHandler{
		RedisBroadCastAdaptor: redisBroadCastAdaptor,
	}
	http.Handle("/api/notify", http.HandlerFunc(notificationHandler.Notify))

	log.Println("Serving at localhost:5007...")
	log.Println("Serving at localhost:5007...")
	log.Fatal(http.ListenAndServe(":5007", nil))
}
