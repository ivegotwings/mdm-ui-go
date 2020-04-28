package notification

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ivegotwings/mdm-ui-go/redis"
)

type notificationInfo struct {
	showNotificationToUser string
	id                     string
	timeStamp              string
	source                 string
	userid                 string
	connectionId           string
}

type ClientState struct {
	notificationInfo notificationInfo
}

type JsonData struct {
	ClientState ClientState `json:"clientState"`
}

type AttributeString struct {
	Locale string `json:"locale"`
	Source string `json:"source"`
	Id     string `json:"id"`
	Value  string `json:"value"`
}

type AttributeInt struct {
	Locale string `json:"locale"`
	Source string `json:"source"`
	Id     string `json:"id"`
	Value  int    `json:"value"`
}

type AttributeStringVal struct {
	Values []AttributeString `json:"values"`
}

type AttributeIntVal struct {
	Values []AttributeInt `json:"values"`
}

type Attributes struct {
	EntityAction     AttributeStringVal `json:"entityAction"`
	EntityId         AttributeStringVal `json:"entityId"`
	EntityType       AttributeStringVal `json:"entityType"`
	RequestId        AttributeStringVal `json:"requestId"`
	RequestStatus    AttributeStringVal `json:"requestStatus"`
	RequestTimestamp AttributeIntVal    `json:"requestTimestamp"`
	RelatedRequestId AttributeStringVal `json:"relatedRequestId"`
	RequestGroupId   AttributeStringVal `json:"requestGroupId"`
	ClientId         AttributeStringVal `json:"clientId"`
	UserId           AttributeStringVal `json:"userId"`
	ObjectStore      AttributeStringVal `json:"ObjectStore"`
}

type Data struct {
	Attributes Attributes `json:"attributes"`
}

type NotificationObject struct {
	Data     Data     `json:"data"`
	JsonData JsonData `json:"jsonData"`
}

type Notification struct {
	NotificationObject NotificationObject `json:"notificationObject"`
}

func Notify(w http.ResponseWriter, r *http.Request, redisBroadCastAdaptor *redis.Broadcast) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}
	var _message Notification
	err = json.Unmarshal(body, &_message)
	if err != nil {
		fmt.Println("ERR", err)
		log.Println("notify error in processing body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		fmt.Printf("/api/notify: %v\n", _message.NotificationObject)
		//fmt.Println("/api/notify message", _message.notificationObject.data.attributes.entityAction)
		redisBroadCastAdaptor.Send(nil, "testroom", "event:notification", _message)
	}
	return nil
}
