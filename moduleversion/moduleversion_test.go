package moduleversion

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestMODULEVERSION(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	LoadDomainMap()
	_moduleDomainMap := GetModuleDomainMap()
	v := reflect.ValueOf(_moduleDomainMap)

	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()
		if len(value.([]string)) == 0 {
			t.Error("ModuleVersion- Load failed")
		} else {
			fmt.Println(len(value.([]string)))
		}
	}
	testDomain := "referenceData"
	resolvedDomain, err := GetResolvedDomain("entityData", testDomain)
	if resolvedDomain != testDomain || err != nil {
		t.Error("ModuleVersion- failed resolveddomain")
	}

	versionkey, error := GetVersionKey("entityData", "referenceData", "rdwengg-az-dev2")
	if error != nil {
		t.Error("ModuleVersion- failed VersionKey")
	}
	if versionkey != "entityData"+"-"+"referenceData"+"-"+"rdwengg-az-dev2"+MODULE_VERSION_KEY {
		t.Error("ModuleVersion- failed VersionKey")
	}
}
