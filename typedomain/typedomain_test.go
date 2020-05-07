package typedomain

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/ivegotwings/mdm-ui-go/executioncontext"
)

func TestTYPEDOMAIN(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	ctx := executioncontext.Context{
		UserId:            "rdwadmin@riversand.com_user",
		TenantId:          "rdwengg-az-dev2",
		ClientId:          "rufClient",
		Host:              "manage.engg-az-dev2.riversand-dataplatform.com:7075",
		UserRoles:         "[\"\"]",
		DefaultRole:       "",
		OwnershipData:     "",
		OwnershipEditData: "",
		UserEmail:         "",
		FirstName:         "",
		LastName:          "",
	}

	domain, err := GetDomainForEntityType("sku", ctx)
	if err != nil {
		t.Error(err)
	} else {
		if domain != "thing" {
			t.Error("typedomain get- incorrect domain")
		}
	}
	domainmap, error := InitializeEntityTypeDomainMap(ctx)
	if error != nil {
		t.Error("typedomain- failed init call " + error.Error())
	}

	if domainmap["sku_entityType"] != "thing" {
		t.Error("typedomain init- incorrect domain")
	}
}
