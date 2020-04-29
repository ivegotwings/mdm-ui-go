package moduleversion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var MODULE_VERSION_KEY string = "-runtime-module-version"

type ModuleDomainMap struct {
	EntityData        []string `json:"entityData"`
	EntityModel       []string `json:"entityModel"`
	EntityGovernData  []string `json:"entityGovernData"`
	Config            []string `json:"config"`
	EventData         []string `json:"eventData"`
	GenericObjectData []string `json:"genericObjectData"`
}

func UpdateModuleVersion(module string, domain string, tenantId string) {
	if domain == "" {
		domain = "default"
	}
}

func LoadDomainMap() {
	var moduleDomainMap ModuleDomainMap
	mapFile, err := os.Open("moduledomainmap.json")
	defer mapFile.Close()
	byteValue, _ := ioutil.ReadAll(mapFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	_ = json.Unmarshal([]byte(byteValue), &moduleDomainMap)
	fmt.Println("LoadDomainMap- ", moduleDomainMap)
}
