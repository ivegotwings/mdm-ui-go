package typedomain

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

var entityTypeDomainLookUp = map[string]string{
	"uomLengthWithoutFormula3_entityType":       "UOMData",
	"uomLengthWithoutFormula2_entityType":       "UOMData",
	"uomLengthMeasure_entityType":               "UOMData",
	"imagerendition_entityType":                 "digitalAsset",
	"warranty_entityType":                       "thing",
	"color_entityType":                          "thing",
	"pallet_entityType":                         "thing",
	"mobile_entityType":                         "thing",
	"sdf_entityType":                            "thing",
	"dfgdfg_entityType":                         "thing",
	"dfgdf_entityType":                          "thing",
	"sdsa_entityType":                           "thing",
	"fgdv_entityType":                           "thing",
	"sizevaluemapping_entityType":               "generic",
	"customer_entityType":                       "party",
	"transportcategory_entityType":              "referenceData",
	"supplier_entityType":                       "referenceData",
	"skuoption1_entityType":                     "referenceData",
	"bottlesizes_entityType":                    "referenceData",
	"packagingtype_entityType":                  "referenceData",
	"productgroup_entityType":                   "thing",
	"bundle_entityType":                         "thing",
	"city_entityType":                           "referenceData",
	"itemtype_entityType":                       "referenceData",
	"colors_entityType":                         "referenceData",
	"warrantytype_entityType":                   "referenceData",
	"entityCompositeModel_entityType":           "sysBaseModel",
	"entityGovernModelModel_entityType":         "sysBaseModel",
	"entityContextModel_entityType":             "sysBaseModel",
	"entityDefaultValuesModel_entityType":       "sysBaseModel",
	"entityType_entityType":                     "baseModel",
	"relationshipModel_entityType":              "baseModel",
	"taxonomy_entityType":                       "taxonomyModel",
	"quantity_entityType":                       "referenceData",
	"weight_entityType":                         "referenceData",
	"volume_entityType":                         "referenceData",
	"profileType_entityType":                    "sysReferenceData",
	"dataType_entityType":                       "sysReferenceData",
	"displayType_entityType":                    "sysReferenceData",
	"eventType_entityType":                      "sysReferenceData",
	"eventSubType_entityType":                   "sysReferenceData",
	"uomFlowrate_entityType":                    "UOMData",
	"uomFocallength_entityType":                 "UOMData",
	"uomImagecapturespeed_entityType":           "UOMData",
	"uomMegapixels_entityType":                  "UOMData",
	"uomRotationspeed_entityType":               "UOMData",
	"uomPageyield_entityType":                   "UOMData",
	"uomAngle_entityType":                       "UOMData",
	"uomVolume_entityType":                      "UOMData",
	"uomTimesmall_entityType":                   "UOMData",
	"uomStoragecapacity_entityType":             "UOMData",
	"uomSensitivity_entityType":                 "UOMData",
	"audio_entityType":                          "digitalAsset",
	"bulkoperationevent_entityType":             "systemEvent",
	"entitymanageevent_entityType":              "systemEvent",
	"contentTemplateModel_entityType":           "governanceModel",
	"systemDashboard_entityType":                "visualization",
	"connectorItemState_entityType":             "sysReferenceData",
	"uomLengthWithoutFormula1_entityType":       "UOMData",
	"uomLengthWithoutFormula_entityType":        "UOMData",
	"connectorListedState_entityType":           "sysReferenceData",
	"connectorIntroState_entityType":            "sysReferenceData",
	"audiorendition_entityType":                 "digitalAsset",
	"organization_entityType":                   "party",
	"each_entityType":                           "thing",
	"kit_entityType":                            "thing",
	"abctype_entityType":                        "thing",
	"style_entityType":                          "thing",
	"dgdfg_entityType":                          "thing",
	"dfg_entityType":                            "thing",
	"zxcxz_entityType":                          "thing",
	"sdas_entityType":                           "thing",
	"fghfg_entityType":                          "thing",
	"colorvaluemapping_entityType":              "generic",
	"connectorchannel_entityType":               "referenceData",
	"ownershipdata_entityType":                  "referenceData",
	"fragrance_entityType":                      "referenceData",
	"skuoption2_entityType":                     "referenceData",
	"accountrecordtype_entityType":              "referenceData",
	"withholdingcodesref_entityType":            "referenceData",
	"enrichitem_entityType":                     "thing",
	"item_entityType":                           "thing",
	"product_entityType":                        "thing",
	"state_entityType":                          "referenceData",
	"addresscountry_entityType":                 "referenceData",
	"role_entityType":                           "sysAuthorizationModel",
	"authorizationModel_entityType":             "sysAuthorizationModel",
	"user_entityType":                           "sysAuthorizationModel",
	"classification_entityType":                 "taxonomyModel",
	"source_entityType":                         "referenceData",
	"percentage_entityType":                     "referenceData",
	"entityDisplayModel_entityType":             "sysBaseModel",
	"entityManageModel_entityType":              "sysBaseModel",
	"attributeModel_entityType":                 "baseModel",
	"interactionLocale_entityType":              "sysReferenceData",
	"area_entityType":                           "referenceData",
	"ruleType_entityType":                       "sysReferenceData",
	"numberFormat_entityType":                   "sysReferenceData",
	"executionMode_entityType":                  "sysReferenceData",
	"relatedRequestId_entityType":               "sysReferenceData",
	"clientId_entityType":                       "sysReferenceData",
	"sortType_entityType":                       "sysReferenceData",
	"uomDatatransferratebits_entityType":        "UOMData",
	"uomDutycycle_entityType":                   "UOMData",
	"uomElectriccurrent_entityType":             "UOMData",
	"uomBatteryaverageruntimecamera_entityType": "UOMData",
	"uomAngularvelocity_entityType":             "UOMData",
	"uomArea_entityType":                        "UOMData",
	"uomPressure_entityType":                    "UOMData",
	"uomLength_entityType":                      "UOMData",
	"variantgenerationevent_entityType":         "systemEvent",
	"entitymanageappevent_entityType":           "systemEvent",
	"entitygovernevent_entityType":              "systemEvent",
	"troubleshootingevent_entityType":           "systemEvent",
	"uomTorque_entityType":                      "UOMData",
	"uomSurgeprotection_entityType":             "UOMData",
	"uomForceperlength_entityType":              "UOMData",
	"uomFrequency_entityType":                   "UOMData",
	"uomMass_entityType":                        "UOMData",
	"uomPrintspeed_entityType":                  "UOMData",
	"uomPower_entityType":                       "UOMData",
	"uomMediacapacity_entityType":               "UOMData",
	"uomMediaquantity_entityType":               "UOMData",
	"uomResolution_entityType":                  "UOMData",
	"workflowDefinition_entityType":             "governanceModel",
	"attributemapping_entityType":               "sysIntegrationModel",
	"contextmapping_entityType":                 "sysIntegrationModel",
	"connectorrequestactivity_entityType":       "generic",
	"image_entityType":                          "digitalAsset",
	"document_entityType":                       "digitalAsset",
	"connectorChannelState_entityType":          "sysReferenceData",
	"videorendition_entityType":                 "digitalAsset",
	"productsku_entityType":                     "thing",
	"electronicsentitytype_entityType":          "thing",
	"case_entityType":                           "thing",
	"testadd_entityType":                        "thing",
	"werwe_entityType":                          "thing",
	"xcxz_entityType":                           "thing",
	"sdfsd_entityType":                          "thing",
	"a_entityType":                              "thing",
	"vendor_entityType":                         "party",
	"employeesizecategory_entityType":           "referenceData",
	"supplierentitytyperef_entityType":          "referenceData",
	"primarysize_entityType":                    "referenceData",
	"integrationchannel_entityType":             "referenceData",
	"supplierownershiptyperef_entityType":       "referenceData",
	"paymentmethodsref_entityType":              "referenceData",
	"polishtype_entityType":                     "referenceData",
	"vendorgroupingkeyref_entityType":           "referenceData",
	"secondarysize_entityType":                  "referenceData",
	"industry_entityType":                       "referenceData",
	"sizes_entityType":                          "referenceData",
	"sku_entityType":                            "thing",
	"country_entityType":                        "referenceData",
	"brand_entityType":                          "referenceData",
	"channel_entityType":                        "referenceData",
	"locale_entityType":                         "referenceData",
	"entityValidationModel_entityType":          "sysBaseModel",
	"length_entityType":                         "referenceData",
	"time_entityType":                           "referenceData",
	"timeFormat_entityType":                     "sysReferenceData",
	"dateFormat_entityType":                     "sysReferenceData",
	"triggerAction_entityType":                  "sysReferenceData",
	"entityAction_entityType":                   "sysReferenceData",
	"activityCriteria_entityType":               "sysReferenceData",
	"entityId_entityType":                       "sysReferenceData",
	"referenceRelationship_entityType":          "sysReferenceData",
	"uomTemperature_entityType":                 "UOMData",
	"uomStorage_entityType":                     "UOMData",
	"uomTime_entityType":                        "UOMData",
	"uomAngleplane_entityType":                  "UOMData",
	"uomBatterycapacity_entityType":             "UOMData",
	"uomForce_entityType":                       "UOMData",
	"uomMassperlength_entityType":               "UOMData",
	"uomElectricalcapacitan_entityType":         "UOMData",
	"uomElectricalinductance_entityType":        "UOMData",
	"uomElectricalresistance_entityType":        "UOMData",
	"uomBrightness_entityType":                  "UOMData",
	"uomDatatransferrate_entityType":            "UOMData",
	"uomDotpitch_entityType":                    "UOMData",
	"uomElectricalpotential_entityType":         "UOMData",
	"uomWeight_entityType":                      "UOMData",
	"workflowDefinitionMapping_entityType":      "governanceModel",
	"connectormessageactivity_entityType":       "sysIntegrationModel",
	"video_entityType":                          "digitalAsset",
	"externalevent_entityType":                  "systemEvent",
	"entitymodelevent_entityType":               "systemEvent",
	"businessCondition_entityType":              "governanceModel",
	"businessRule_entityType":                   "governanceModel",
	"graphProcessModel_entityType":              "governanceModel",
	"ruleContextMappings_entityType":            "governanceModel",
	"healthcheck_entityType":                    "generic",
	"connectorValidationState_entityType":       "sysReferenceData",
	"connectorSyndicationState_entityType":      "sysReferenceData",
	"dashboard_entityType":                      "visualization",
}

type EntityModel struct {
	Domain string `json:"domain"`
}

type TypeDomainResponseBody struct {
	EntityModels []EntityModel `json:"entityModels"`
}

type TypeDomainResponse struct {
	Response TypeDomainResponseBody `json:"response"`
}

// {
// 	"request":{
// 	   "returnRequest":false,
// 	   "params":{

// 	   },
// 	   "requestId":"701e1dde-5505-424f-a445-b922bf8ea57d"
// 	},
// 	"response":{
// 	   "entityModels":[
// 		  {
// 			 "id":"uomLengthWithoutFormula3_entityType",
// 			 "name":"uomLengthWithoutFormula3",
// 			 "type":"entityType",
// 			 "domain":"UOMData",
// 			 "source":"internal",
// 			 "properties":{
// 				"externalName":"Length Without Formula3",
// 				"baseUnitSymbol":"m",
// 				"createdService":"entityManageModelService",
// 				"createdBy":"rdwadmin@riversand.com_user",
// 				"modifiedService":"entityManageModelService",
// 				"modifiedBy":"rdwadmin@riversand.com_user",
// 				"createdDate":"2020-01-04T11:00:57.532-0600",
// 				"modifiedDate":"2020-01-04T11:00:57.532-0600"
// 			 },
// 			 "data":{
// 				"attributes":{
// 				   "baseUnitSymbol":{
// 					  "values":[
// 						 {
// 							"locale":"en-US",
// 							"source":"internal",
// 							"id":"70d9e38e-91df-4244-bed6-d0c4dc7f86b9",
// 							"value":"m"
// 						 }
// 					  ]
// 				   },
// 				   "externalName":{
// 					  "values":[
// 						 {
// 							"locale":"en-US",
// 							"source":"internal",
// 							"id":"4a64d1a6-e5d7-489e-84fa-3c91fd8ed122",
// 							"value":"Length Without Formula3"
// 						 }
// 					  ]
// 				   }
// 				}
// 			 }
// 		  }
// 	   ],
// 	   "status":"success",
// 	   "totalRecords":1
// 	}
//  }

func GetDomainForEntityType(entityType string) (string, error) {
	lookUpValue := entityTypeDomainLookUp[entityType+"_entityType"]
	if lookUpValue == "" {
		//post call
		var requestBody []byte = []byte(`{"params":{"query":{"ids":["` + entityType + `_entityType"],"filters":{"typesCriterion":["entityType"]}},"fields": {"attributes": ["_ALL"],"relationships": ["_ALL"]}}}`)
		req, err := http.NewRequest("POST", "http://manage.engg-az-dev2.riversand-dataplatform.com:8085/rdwengg-az-dev2/api/entitymodelservice/get", bytes.NewBuffer(requestBody))
		if err != nil {
			return "", err
		} else {
			req.Header.Set("x-rdp-tenantId", "rdwengg-az-dev2")
			req.Header.Set("x-rdp-userId", "rdwadmin@riversand.com_user")
			req.Header.Set("x-rdp-userRoles", "[\"admin\"]")
			req.Header.Set("x-rdp-useremail", "rdwadmin@riversand.com")
			req.Header.Set("x-rdp-defaultrole", "admin")
			req.Header.Set("x-rdp-clientid", "rufClient")
			req.Header.Set("x-rdp-ownershipdata", "")
			req.Header.Set("x-rdp-ownershipeditdata", "")
			req.Header.Set("x-rdp-useremail", "rdwadmin@riversand.com")
			req.Header.Set("x-rdp-firstName", "rdw")
			req.Header.Set("x-rdp-lastName", "admin")

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("x-rdp-authtoken", "m4eZW93FLaUAUfoR1vYEEfwTXr1wdbedZNss0aId6CQ=")

			client := &http.Client{
				Timeout: 30 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				return "", err
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return "", err
				}
				var response TypeDomainResponse
				if json.Unmarshal(body, &response) != nil {
					return "", err
				}
				lookUpValue = response.Response.EntityModels[0].Domain
				if lookUpValue == "" {
					return "", errors.New("doamin not found")
				}
			}
			defer resp.Body.Close()
		}
	}
	return lookUpValue, nil
}
