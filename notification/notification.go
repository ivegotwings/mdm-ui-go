package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ivegotwings/mdm-ui-go/redis"
	"github.com/ivegotwings/mdm-ui-go/utils"
)

type UserNotificationInfo struct {
	NotificationInfo
	RequestStatus string
	TaskId        string
	TaskType      string
	RequestId     string
	ServiceName   string
	Description   string
	Status        string
}

type Context struct {
	AppInstanceId string `json:"appInstanceId"`
	Id            string `json:"id"`
	Type          string `json:"type"`
	DataIndex     string `json:"dataIndex"`
}

type NotificationInfo struct {
	ShowNotificationToUser bool    `json:"showNotificationToUser"`
	Id                     string  `json:"id"`
	TimeStamp              string  `json:"timeStamp"`
	Source                 string  `json:"source"`
	UserId                 string  `json:"userId"`
	ConnectionId           string  `json:"connectionId"`
	Context                Context `json:"context"`
	ActionType             string  `json:"actionType"`
	Operation              string  `json:"operation"`
}

type ClientState struct {
	NotificationInfo     NotificationInfo
	EmulatedSyncDownload bool
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
	ServiceName      AttributeStringVal `json:"serviceName"`
	TaskId           AttributeStringVal `json:"taskId"`
	TaskType         AttributeStringVal `json:"taskType"`
}

type Data struct {
	Attributes Attributes `json:"attributes"`
	JsonData   JsonData   `json:"jsonData"`
}

type NotificationObject struct {
	Data Data `json:"data"`
}

type Notification struct {
	NotificationObject NotificationObject `json:"notificationObject"`
	TenantId           string             `json:"tenantId"`
	ServiceName        string             `json:"serviceName"`
	Domain             string             `json:"domain"`
	Params             interface{}        `json:"params"`
	ReturnRequest      bool               `json:"returnRequest"`
	Id                 string             `json:"id"`
	Type               string             `json:"type"`
	Properties         Properties         `json:"properties"`
}

type Properties struct {
	CreatedService  string `json:"createdService"`
	CreatedBy       string `json:"createdBy"`
	ModifiedService string `json:"modifiedService"`
	ModifiedBy      string `json:"modifiedBy"`
	CreatedDate     string `json:"createdDate"`
	ModifiedDate    string `json:"modifiedDate"`
}

var clientIdNotificationExlusionList = []string{"healthcheckClient"}

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
		return err
	} else {
		fmt.Printf("Notify: %v\n", _message.NotificationObject)
		tenantId := _message.TenantId
		userId := _message.NotificationObject.Data.JsonData.ClientState.NotificationInfo.UserId
		if tenantId != "" && userId != "" {
			var clientId string
			if len(_message.NotificationObject.Data.Attributes.ClientId.Values) > 0 {
				clientId = _message.NotificationObject.Data.Attributes.ClientId.Values[0].Value
			}
			if clientId != "" {
				if ok := utils.Contains(clientIdNotificationExlusionList, clientId); ok {
					fmt.Println("Ignoring notification for clientId", clientId)
				}
				return sendNotification(_message.NotificationObject, tenantId)
			} else {
				err = errors.New("Notify- missing clientId")
				return err
			}
		} else {
			err = errors.New("Notify- tenantId or userId not found")
			return err
		}
	}
}

func sendNotification(notificationObject NotificationObject, tenantId string) error {
	//redisBroadCastAdaptor.Send(nil, "testroom", "event:notification", _message)
	var userNotificationInfo UserNotificationInfo
	err := prepareNotificationObject(&userNotificationInfo, notificationObject)
	if err != nil {
		fmt.Println("sendNotification- error in pepareNotificationObject ", err)
		return err
	} else {
		fmt.Printf("sendNotication userNotificationInfo: %v\n", userNotificationInfo)
	}
	return nil
}

func prepareNotificationObject(userNotificationInfo *UserNotificationInfo, notificationObject NotificationObject) error {
	var entityId, entityType string
	var err error
	if len(notificationObject.Data.Attributes.EntityId.Values) > 0 {
		entityId = notificationObject.Data.Attributes.EntityId.Values[0].Value
	}

	if len(notificationObject.Data.Attributes.EntityType.Values) > 0 {
		entityType = notificationObject.Data.Attributes.EntityType.Values[0].Value
	}

	if entityId == "" || entityType == "" {
		err = errors.New("prepareNotificationObject- missing entityId or entityType")
	} else {
		userNotificationInfo.ShowNotificationToUser = notificationObject.Data.JsonData.ClientState.NotificationInfo.ShowNotificationToUser
		userNotificationInfo.Id = notificationObject.Data.JsonData.ClientState.NotificationInfo.Id
		userNotificationInfo.TimeStamp = notificationObject.Data.JsonData.ClientState.NotificationInfo.TimeStamp
		userNotificationInfo.Source = notificationObject.Data.JsonData.ClientState.NotificationInfo.Source
		userNotificationInfo.UserId = notificationObject.Data.JsonData.ClientState.NotificationInfo.UserId
		userNotificationInfo.ConnectionId = notificationObject.Data.JsonData.ClientState.NotificationInfo.ConnectionId
		userNotificationInfo.Context = notificationObject.Data.JsonData.ClientState.NotificationInfo.Context
		if len(notificationObject.Data.Attributes.RequestStatus.Values) > 0 {
			userNotificationInfo.RequestStatus = notificationObject.Data.Attributes.RequestStatus.Values[0].Value
		}
	}

	var desc string

	if userNotificationInfo.Context.Id != entityId {
		userNotificationInfo.ShowNotificationToUser = false
		desc = "System Manage Complete"
	}
	userNotificationInfo.Context.Id = entityId
	userNotificationInfo.Context.Type = entityType
	if len(notificationObject.Data.Attributes.ServiceName.Values) > 0 {
		userNotificationInfo.ServiceName = notificationObject.Data.Attributes.ServiceName.Values[0].Value
	}
	if len(notificationObject.Data.Attributes.TaskId.Values) > 0 {
		userNotificationInfo.TaskId = notificationObject.Data.Attributes.TaskId.Values[0].Value
	}
	if len(notificationObject.Data.Attributes.TaskType.Values) > 0 {
		userNotificationInfo.TaskType = notificationObject.Data.Attributes.TaskType.Values[0].Value
	}

	if userNotificationInfo.Operation == "" {
		userNotificationInfo.Operation = "connectIntegrationType"
	}

	switch status := strings.ToLower(userNotificationInfo.RequestStatus); status {
	case "completed":
		userNotificationInfo.Status = "success"
		break
	case "errored":
		userNotificationInfo.Status = "error"
		break
	default:
		userNotificationInfo.Status = userNotificationInfo.RequestStatus
	}

	userNotificationInfo.Description = desc
	return err
}
