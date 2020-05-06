package executioncontext

import (
	"net/http"
	"os"
)

type Context struct {
	UserId            string
	TenantId          string
	ClientAuthKey     string
	ClientId          string
	OwnershipData     string
	OwnershipEditData string
	FirstName         string
	LastName          string
	FullName          string
	UserName          string
	UserEmail         string
	UserRoles         string
	DefaultRole       string
	Referer           string
	Host              string
}

// 'clientAuthKey': clientAuthKey ? clientAuthKey : "",
// 'headers': {
// 	"clientId": clientId ? clientId : "",
// 	"ownershipData": req.headers["x-rdp-ownershipdata"] ? JSON.parse(req.headers["x-rdp-ownershipdata"]) : userDefaults.ownershipData,
// 	"ownershipEditData": req.headers["x-rdp-ownershipeditdata"] ? JSON.parse(req.headers["x-rdp-ownershipeditdata"]) : userDefaults.ownershipEditData,
// 	"userId": uid,
// 	"firstName": firstName,
// 	"lastName": lastName,
// 	"fullName": fullName,
// 	"userName": req.headers["x-rdp-username"] || userDefaults.userName,
// 	"userEmail": req.headers["x-rdp-useremail"] || userDefaults.userEmail,
// 	"userRoles": roles,
// 	"defaultRole": defaultRole,
// 	"isTrustBasedOnUserModel": isTrustBasedOnUserModel

func GetContext(req *http.Request) Context {
	UserContext := Context{
		UserId:            req.Header.Get("x-rdp-userid"),
		TenantId:          req.Header.Get("x-rdp-tenantid"),
		UserRoles:         req.Header.Get("x-rdp-userroles"),
		DefaultRole:       req.Header.Get("x-rdp-defaultrole"),
		UserName:          req.Header.Get("x-rdp-username"),
		UserEmail:         req.Header.Get("x-rdp-useremail"),
		FirstName:         req.Header.Get("x-rdp-firstname"),
		LastName:          req.Header.Get("x-rdp-lastname"),
		OwnershipData:     req.Header.Get("x-rdp-ownershipdata"),
		OwnershipEditData: req.Header.Get("x-rdp-ownershipeditdata"),
		ClientId:          "rufClient",
		ClientAuthKey:     "3218fa37-f809-4be4-b88e-653419b20e28",
		Host:              "rdp-rest:7075",
	}

	if os.Getenv("ENV") == "DEVELOPMENT" {
		UserContext.Host = "manage.engg-az-dev2.riversand-dataplatform.com:7075"
	}
	return UserContext
}
