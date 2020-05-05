package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/ivegotwings/mdm-ui-go/connection"
	"github.com/ivegotwings/mdm-ui-go/moduleversion"
	"github.com/ivegotwings/mdm-ui-go/state"
	"github.com/ivegotwings/mdm-ui-go/typedomain"
	"github.com/ivegotwings/mdm-ui-go/utils"
)

type UserNotificationInfo struct {
	NotificationInfo
	RequestStatus        string
	TaskId               string
	TaskType             string
	RequestId            string
	ServiceName          string
	Description          string
	Status               string
	Action               int
	DataIndex            string
	EmulatedSyncDownload bool
	TenantId             string
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
	EntityAction           AttributeStringVal `json:"entityAction"`
	EntityId               AttributeStringVal `json:"entityId"`
	EntityType             AttributeStringVal `json:"entityType"`
	RequestId              AttributeStringVal `json:"requestId"`
	RequestStatus          AttributeStringVal `json:"requestStatus"`
	RequestTimestamp       AttributeIntVal    `json:"requestTimestamp"`
	RelatedRequestId       AttributeStringVal `json:"relatedRequestId"`
	RequestGroupId         AttributeStringVal `json:"requestGroupId"`
	ClientId               AttributeStringVal `json:"clientId"`
	UserId                 AttributeStringVal `json:"userId"`
	ObjectStore            AttributeStringVal `json:"ObjectStore"`
	ServiceName            AttributeStringVal `json:"serviceName"`
	TaskId                 AttributeStringVal `json:"taskId"`
	TaskType               AttributeStringVal `json:"taskType"`
	ConnectIntegrationType AttributeStringVal `json:"connectIntegrationType"`
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

var actionLookUpTable = map[string]string{
	"MODEL_IMPORT_success":                                                        "ModelImportComplete",
	"MODEL_IMPORT_success_but_errors":                                             "ModelImportCompletedWithErrors",
	"MODEL_IMPORT_error":                                                          "ModelImportFail",
	"MODEL_EXPORT_success_true":                                                   "EmulatedSyncDownloadComplete",
	"MODEL_EXPORT_success_false":                                                  "RSConnectComplete",
	"MODEL_EXPORT_error_":                                                         "RSConnectFail",
	"MODEL_EXPORT_success_but_errors_":                                            "RSConnectFail",
	"ENTITY_EXPORT_success_true":                                                  "EmulatedSyncDownloadComplete",
	"ENTITY_EXPORT_success_false":                                                 "RSConnectComplete",
	"ENTITY_EXPORT_error_false":                                                   "RSConnectFail",
	"ENTITY_EXPORT_success_but_errors_false":                                      "RSConnectFail",
	"configurationmanageservice_uiconfig_success":                                 "ConfigurationSaveComplete",
	"configurationmanageservice_uiconfig_error":                                   "ConfigurationSaveFail",
	"entitymanageservice_success_System.Manage.Complete":                          "SystemSaveComplete",
	"entitymanageservice_success_default":                                         "SaveComplete",
	"entitymanageservice_error_System.Manage.Complete":                            "SystemSaveFail",
	"entitymanageservice_error_default":                                           "SaveFail",
	"entitymanagemodelservice_sucess_default":                                     "ModelSaveComplete",
	"entitymanagemodelservice_error_default":                                      "ModelSaveFail",
	"entitymanagemodelservice_success":                                            "ModelSaveComplete",
	"entitymanagemodelservice_error":                                              "ModelSaveFail",
	"entitygovernservice_WorkflowTransition_success":                              "WorkflowTransitionComplete",
	"entitygovernservice_WorkflowTransition_error":                                "WorkflowTransitionFail",
	"entitygovernservice_WorkflowAssignment_success":                              "WorkflowAssignmentComplete",
	"entitygovernservice_WorkflowAssignment_error":                                "WorkflowAssignmentFail",
	"entitygovernservice_BusinessCondition_success":                               "BusinessConditionSaveComplete",
	"entitygovernservice_BusinessCondition_error":                                 "BusinessConditionSaveFail",
	"entitygovernservice_success":                                                 "GovernComplete",
	"entitygovernservice_error":                                                   "GovernFail",
	"notificationmanageservice_changeAssignment-multi-query_success":              "BulkWorkflowAssignmentComplete",
	"notificationmanageservice_changeAssignment-multi-query_success_but_errors":   "BulkWorkflowAssignmentComplete",
	"notificationmanageservice_changeAssignment-multi-query_error":                "WorkflowAssignmentFail",
	"notificationmanageservice_transitionWorkflow-multi-query_success":            "BulkWorkflowTransitionComplete",
	"notificationmanageservice_transitionWorkflow-multi-query_success_but_errors": "BulkWorkflowTransitionComplete",
	"notificationmanageservice_transitionWorkflow-multi-query_error":              "WorkflowTransitionFail",
}

var actions = map[string]int{
	"SystemSaveComplete":             1,
	"SaveComplete":                   2,
	"SystemSaveFail":                 3,
	"SaveFail":                       4,
	"GovernComplete":                 5,
	"GovernFail":                     6,
	"WorkflowTransitionComplete":     7,
	"WorkflowTransitionFail":         8,
	"WorkflowAssignmentComplete":     9,
	"BulkWorkflowAssignmentComplete": 20,
	"BulkWorkflowTransitionComplete": 21,
	"WorkflowAssignmentFail":         10,
	"RSConnectComplete":              11,
	"RSConnectFail":                  12,
	"BusinessConditionSaveComplete":  13,
	"BusinessConditionSaveFail":      14,
	"ModelImportComplete":            15,
	"ModelImportFail":                16,
	"ModelSaveComplete":              17,
	"ModelSaveFail":                  18,
	"EmulatedSyncDownloadComplete":   19,
	"ConfigurationSaveComplete":      22,
	"ConfigurationSaveFail":          23,
	"ModelImportCompletedWithErrors": 24,
}

var dataIndexMapping = map[string]string{
	"entitymanageService":        "entityData",
	"entitygovernservice":        "entityData",
	"entitymanagemodelservice":   "entityModel",
	"configurationmanageservice": "config",
	"genericobjectmanageservice": "genericObjectData",
}

var clientIdNotificationExlusionList = []string{"healthcheckClient"}

var redisBroadCastAdaptor *connection.Broadcast

type NotificationHandler struct {
	RedisBroadCastAdaptor connection.Broadcast
}

type NotificationPayload struct {
	UserNotificationInfo UserNotificationInfo
	UserInfo             map[string]string
	VersionKey           string
}

var NotificationPayloadChannel = make(chan NotificationPayload)

func SetRedisBroadCastAdaptor(adaptor *connection.Broadcast) {
	redisBroadCastAdaptor = adaptor
	var quit = make(chan struct{})
	go NotificationScheduler(quit)

}

func (notificationHandler *NotificationHandler) Notify(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		go processNotification(w, r, notificationHandler)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	fmt.Fprint(w, "POST done")
	w.Header().Set("Server", "A Go Web Server")
	w.WriteHeader(200)

}

func processNotification(w http.ResponseWriter, r *http.Request, notificationHandler *NotificationHandler) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}
	var _message Notification
	err = json.Unmarshal(body, &_message)
	if err != nil {
		utils.PrintInfo("notify error in processing body: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	} else {
		log.Printf("Notify: %v\n", _message.NotificationObject)
		tenantId := _message.TenantId
		userId := _message.NotificationObject.Data.JsonData.ClientState.NotificationInfo.UserId
		if tenantId != "" && userId != "" {
			var clientId string
			if len(_message.NotificationObject.Data.Attributes.ClientId.Values) > 0 {
				clientId = _message.NotificationObject.Data.Attributes.ClientId.Values[0].Value
			}
			if clientId != "" {
				if ok := utils.Contains(clientIdNotificationExlusionList, clientId); ok {
					utils.PrintInfo("Ignoring notification for clientId: " + clientId)
				}
				go sendNotification(_message.NotificationObject, tenantId)
			} else {
				err = errors.New("Notify- missing clientId")
				return err
			}
		} else {
			err = errors.New("Notify- tenantId or userId not found")
			return err
		}
	}
	return nil
}

func sendNotification(notificationObject NotificationObject, tenantId string) error {
	var userNotificationInfo UserNotificationInfo
	err := prepareNotificationObject(&userNotificationInfo, notificationObject)
	if err != nil {
		utils.PrintInfo("sendNotification- error in pepareNotificationObject " + err.Error())
		return err
	} else {
		if userNotificationInfo.UserId == "" && userNotificationInfo.RequestStatus == "error" {
			return errors.New("sendNotification- Invalid userId or RequestStatus")
		}
		userNotificationInfo.TenantId = tenantId
		userInfo := map[string]string{
			"userId":   userNotificationInfo.UserId,
			"tenantId": userNotificationInfo.TenantId,
		}
		if userNotificationInfo.Action == actions["ModelImportComplete"] || userNotificationInfo.Action == actions["ModelImportCompletedWithErrors"] {
			versionKey, error := moduleversion.GetVersionKey(userNotificationInfo.DataIndex, "", tenantId)
			if error == nil {
				NotificationPayloadChannel <- NotificationPayload{
					VersionKey:           versionKey,
					UserNotificationInfo: userNotificationInfo,
					UserInfo:             userInfo,
				}
			} else {
				return errors.New("sendNotification- cannot get VersionKey")
			}
			//NotificaitonPayloadQueue.Payload
		} else {
			typeDomain, err := typedomain.GetDomainForEntityType(userNotificationInfo.Context.Type)
			if err != nil {
				return err
			}
			versionKey, error := moduleversion.GetVersionKey(userNotificationInfo.DataIndex, typeDomain, tenantId)
			if error == nil {
				NotificationPayloadChannel <- NotificationPayload{
					VersionKey:           versionKey,
					UserNotificationInfo: userNotificationInfo,
					UserInfo:             userInfo,
				}
				utils.PrintInfo("adding payload")
			} else {
				return errors.New("sendNotification- cannot get VersionKey")
			}
		}
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
		//fill userNotificationInfo
		userNotificationInfo.ShowNotificationToUser = notificationObject.Data.JsonData.ClientState.NotificationInfo.ShowNotificationToUser
		userNotificationInfo.Id = notificationObject.Data.JsonData.ClientState.NotificationInfo.Id
		userNotificationInfo.TimeStamp = notificationObject.Data.JsonData.ClientState.NotificationInfo.TimeStamp
		userNotificationInfo.Source = notificationObject.Data.JsonData.ClientState.NotificationInfo.Source
		userNotificationInfo.UserId = notificationObject.Data.JsonData.ClientState.NotificationInfo.UserId
		userNotificationInfo.ConnectionId = notificationObject.Data.JsonData.ClientState.NotificationInfo.ConnectionId
		userNotificationInfo.Context = notificationObject.Data.JsonData.ClientState.NotificationInfo.Context
		userNotificationInfo.EmulatedSyncDownload = notificationObject.Data.JsonData.ClientState.EmulatedSyncDownload
		userNotificationInfo.Operation = notificationObject.Data.JsonData.ClientState.NotificationInfo.Operation
		if userNotificationInfo.Context.Id == "" {
			userNotificationInfo.Context.Id = entityId
			userNotificationInfo.Context.Type = entityType
		}
		if userNotificationInfo.Operation == "" {
			if len(notificationObject.Data.Attributes.ConnectIntegrationType.Values) > 0 {
				userNotificationInfo.Operation = notificationObject.Data.Attributes.ConnectIntegrationType.Values[0].Value
			}
		}
		if len(notificationObject.Data.Attributes.RequestStatus.Values) > 0 {
			userNotificationInfo.RequestStatus = notificationObject.Data.Attributes.RequestStatus.Values[0].Value
		}
	}
	if len(notificationObject.Data.Attributes.ServiceName.Values) > 0 {
		userNotificationInfo.ServiceName = strings.ToLower(notificationObject.Data.Attributes.ServiceName.Values[0].Value)
	}
	if len(notificationObject.Data.Attributes.TaskId.Values) > 0 {
		userNotificationInfo.TaskId = notificationObject.Data.Attributes.TaskId.Values[0].Value
	}
	if len(notificationObject.Data.Attributes.TaskType.Values) > 0 {
		userNotificationInfo.TaskType = notificationObject.Data.Attributes.TaskType.Values[0].Value
	}
	switch status := strings.ToLower(userNotificationInfo.RequestStatus); status {
	case "completed":
		userNotificationInfo.Status = "success"
		break
	case "completed with errors":
		userNotificationInfo.Status = "success_but_errors"
	case "errored":
		userNotificationInfo.Status = "error"
		break
	default:
		userNotificationInfo.Status = strings.ToLower(userNotificationInfo.RequestStatus)
	}

	var desc string = "default"
	if userNotificationInfo.Context.Id != entityId {
		userNotificationInfo.ShowNotificationToUser = false
		desc = "System.Manage.Complete"
	}
	userNotificationInfo.Description = desc

	action, dataIndex := 0, "default"
	if userNotificationInfo.Operation == "MODEL_IMPORT" {
		dataIndex = "entityModel"
		action = actions[actionLookUpTable[userNotificationInfo.Operation+"_"+userNotificationInfo.Status]]
	} else if userNotificationInfo.Operation == "MODEL_EXPORT" || userNotificationInfo.Operation == "ENTITY_EXPORT" {
		action = actions[actionLookUpTable[userNotificationInfo.Operation+"_"+userNotificationInfo.Status+"_"+strconv.FormatBool(userNotificationInfo.EmulatedSyncDownload)]]
	} else if userNotificationInfo.ServiceName == "configurationmanageservice" && userNotificationInfo.Context.Type == "uiconfig" {
		if userNotificationInfo.Operation == "" {
			action = actions[actionLookUpTable[userNotificationInfo.ServiceName+"_"+userNotificationInfo.Context.Type+"_"+userNotificationInfo.Status]]
		}
	} else if userNotificationInfo.ServiceName == "entitymanageservice" {
		if userNotificationInfo.Operation == "" {
			action = actions[actionLookUpTable[userNotificationInfo.ServiceName+"_"+userNotificationInfo.Status+"_"+userNotificationInfo.Description]]
		}
	} else if userNotificationInfo.ServiceName == "entitymanagemodelservice" {
		action = actions[actionLookUpTable[userNotificationInfo.ServiceName+"_"+userNotificationInfo.Status]]
	} else if userNotificationInfo.ServiceName == "entitygovernservice" {
		if userNotificationInfo.Operation == "" {
			action = actions[actionLookUpTable[userNotificationInfo.ServiceName+"_"+userNotificationInfo.Status]]
		} else {
			action = actions[actionLookUpTable[userNotificationInfo.ServiceName+"_"+userNotificationInfo.Operation+"_"+userNotificationInfo.Status]]
		}
	} else if userNotificationInfo.ServiceName == "notificationmanageservice" {
		action = actions[actionLookUpTable[userNotificationInfo.ServiceName+"_"+userNotificationInfo.TaskType+"_"+userNotificationInfo.Status]]
	}
	if val, ok := dataIndexMapping[userNotificationInfo.ServiceName]; ok {
		dataIndex = val
	}
	//log.Printf("setActionAndDataIndex- action & dataIndex", action, dataIndex)

	userNotificationInfo.Action = action
	userNotificationInfo.DataIndex = dataIndex
	return err
}

func NotificationScheduler(quit chan struct{}) {
	for {
		select {
		case payload := <-NotificationPayloadChannel:
			var uniqueVersionKeys = map[string]string{}
			if uniqueVersionKeys[payload.VersionKey] != "done" {
				uniqueVersionKeys[payload.VersionKey] = "done"
				conn := *state.Conn()
				version, err := redis.Int(conn.Do("GET", payload.VersionKey))
				//conn.Flush()
				//version, err := conn.Receive()
				if err == nil {
					var newversion uint8
					//log.Printf("MotificationScheduler versionKey verion", payload.VersionKey, version)
					if version != 0 {
						_version := uint8(version)
						newversion = _version + 1
					} else {
						newversion = moduleversion.DEFAULT_VERSION
					}
					conn.Send("SET", payload.VersionKey, newversion)
					conn.Flush()
				}
			}
			var room string
			if payload.UserInfo["tenantId"] != "" && payload.UserInfo["userId"] != "" {
				room = "socket_conn_room_tenant_" + payload.UserInfo["tenantId"] + "_user_" + payload.UserInfo["userId"]
				//utils.Printf("Broadcasting to room:" + room)
				redisBroadCastAdaptor.Send(nil, room, "event:notification", payload.UserNotificationInfo)
			} else if payload.UserInfo["tenantId"] != "" {
				room = "socket_conn_room_tenant_" + payload.UserInfo["tenantId"]
				redisBroadCastAdaptor.Send(nil, room, "event:notification", payload.UserNotificationInfo)
			}
		case <-quit:
			return
		}
	}
}
