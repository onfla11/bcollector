// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package bconnector // import "github.com/open-telemetry/opentelemetry-collector-contrib/connector/bconnector"

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"reflect"
	"sort"
	"strconv"
	"time"

	"go.uber.org/zap"
	"log"
	"strings"
)

type eventGenerator struct{}

// consume takes a single trace and generate biq event
func (e *eventGenerator) generate(logger *zap.Logger, trace ptrace.Traces) (*plog.LogRecord, error) {
	traceID, err := getTraceID(trace)
	if err != nil {
		return nil, err
	}
	event := plog.NewLogRecord()
	logger.Info("Event generate function called: tenant id: ", zap.Any("trace-id", traceID))
	//body := event.Body()
	//resourceSpans := trace.ResourceSpans()
	//for i := 0; i < resourceSpans.Len(); i++ {
	//	rs := resourceSpans.At(i)
	//	attr := rs.Resource().Attributes()
	//	fmt.Println(attr)
	//}
	//attrMap := body.Map()
	//attrMap.PutStr("test1", "val1")
	//attrMap.PutStr("test2", "val2")
	//attrMap.PutStr("test3", "val3")
	//attrMap.PutStr("traceId", traceID.String())
	event.SetSeverityNumber(plog.SeverityNumberFatal)
	event.SetSeverityText("FATAL")
	event.SetTraceID(traceID)
	logger.Info("Event generated:", zap.Any("event", event))
	return &event, nil
}

func (e *eventGenerator) generateLog(logger *zap.Logger, trace ptrace.Traces) (*plog.Logs, error) {
	var resourceAttr pcommon.Map
	var spanAttr pcommon.Map
	var spanId string
	for i := 0; i < trace.ResourceSpans().Len(); i++ {
		resourceSpans := trace.ResourceSpans().At(i)
		resourceAttr = resourceSpans.Resource().Attributes()
		logger.Info("resourceAttr:", zap.Any("event", resourceAttr))
		scopeSpans := resourceSpans.ScopeSpans()
		for j := 0; j < scopeSpans.Len(); j++ {
			scope := scopeSpans.At(j)
			spans := scope.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				logger.Info("span:", zap.Any("span", span))
				spanId = span.SpanID().String()
				spanAttr = span.Attributes()
			}
		}
	}
	traceID, err := getTraceID(trace)
	if err != nil {
		return nil, err
	}
	logger.Info("Event generateLog function called: tenant id: ", zap.Any("trace-id", traceID))

	//eType := entityMetadataType{Name: "business_transaction", Namespace: namespace{Name: "apm", Version: 1}}

	id := generateId(logger, spanAttr)
	logger.Info("BT id generated: ", zap.Any("bt id", id))
	entity_id := fmt.Sprintf("apm:business_transactio:%s", id)

	btEvent := &bEvent{
		//Data:    `{"kvlistValue":{"values":[{"key": "name","value": {"stringValue": "TEST"}}]}}`,
		Data:    "biq-demo",
		TraceID: traceID.String(),
		SpanID:  spanId,
		//EntityMetadata: entityMetadata{eType},
		Name:         "biq_test",
		Domain:       "biqconnector",
		Bt:           id,
		IsEvent:      true,
		EventType:    "biqconnector:biq_test",
		AppEventType: "biqconnector:biq_test",
		EntityId:     entity_id,
	}

	str := fmt.Sprintf("%+v", btEvent)
	logger.Info("Event generated string val:", zap.Any("eventStr", str))

	attributes := pcommon.NewMap()
	attributes.PutStr("stringValue", str)
	//attributes.PutStr("anotherKey", "anotherValue")

	//myMap := pcommon.NewMap()
	//myMap.PutStr("eventType", "test1233")
	//myMap.PutStr("eventType2", "aaaa1233")

	value, exist := attributes.Get("stringValue")

	logger.Info("Event attributes:", zap.Any("attributes present: ", exist))
	logger.Info("Event attributes:", zap.Any("attributes value: ", value.AsString()))
	logger.Info("Event attributes:", zap.Any("attributes value: ", attributes.AsRaw()))
	jsonStr := `{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"testHost"}}, {"key":"service.namespace","value":{"stringValue":"biqconnector"}}, {"key":"bt.name","value":{"stringValue":"biqPreviewCreateConfig"}},  {"key":"telemetry.sdk.language","value":{"stringValue":"java"}},  {"key":"telemetry.sdk.name","value":{"stringValue":"biqconnector"}},  {"key":"telemetry.sdk.version","value":{"stringValue":"1.11.13"}}],"droppedAttributesCount":0},"scopeLogs":[{"scope":{"name":"name","version":"version","droppedAttributesCount":0},"logRecords":[{"timeUnixNano":"1684617382541971000","observedTimeUnixNano":"1684623646539558000","severityNumber":17,"severityText":"Error","body":{"stringValue":"hello world"},"attributes":[{"key":"sdkVersion","value":{"stringValue":"1.0.1"}}],"droppedAttributesCount":0,"flags":1,"traceId":"0102030405060708090a0b0c0d0e0f10","spanId":"1112131415161718"}],"schemaUrl":"scope_schema"}],"schemaUrl":"resource_schema"}]}`

	//replacementAttributes := []map[string]interface{}{
	//	{
	//		"key": "eventType",
	//		"value": map[string]interface{}{
	//			"stringValue": btEvent.Data,
	//		},
	//	},
	//	{
	//		"key": "EntityMetadata",
	//		"value": map[string]interface{}{
	//			"stringValue": btEvent.EntityMetadata.Data.Name,
	//		},
	//	},
	//	{
	//		"key": "traceId",
	//		"value": map[string]interface{}{
	//			"stringValue": btEvent.TraceID,
	//		},
	//	},
	//}

	logger.Info("resourceAttr: ", zap.Any("resourceAttr", resourceAttr))

	replacementAttributes := make(map[string]interface{})

	for key, value := range resourceAttr.AsRaw() {
		attribute := map[string]interface{}{
			"key": key,
			"value": map[string]interface{}{
				"stringValue": value,
			},
		}
		replacementAttributes[key] = attribute
	}

	val := reflect.ValueOf(*btEvent)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name
		// Get the JSON tag value from the struct field
		jsonTag := val.Type().Field(i).Tag.Get("json")

		// Use the JSON tag value if it exists, otherwise use the struct field name
		if jsonTag != "" {
			fieldName = jsonTag
		}
		fieldValue, _ := strconv.Unquote(ToJSON(field.Interface()))
		if fieldValue == "" {
			fieldValue = ToJSON(field.Interface())
		}

		attribute := map[string]interface{}{
			"key": fieldName,
			"value": map[string]interface{}{
				"stringValue": fieldValue,
			},
		}

		replacementAttributes[fieldName] = attribute
	}

	logger.Info("replacementAttributes: ", zap.Any("replacementAttributes", replacementAttributes))

	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
		return nil, err
	}

	//bodyValue := pcommon.NewValueMap().SetEmptyMap()
	//
	//bodyValue.CopyTo(myMap)
	//for key, value := range myMap {
	//	bodyValue.SetStr()
	//	bodyValue.PutStr(key, value.(string))
	//}
	//data["resourceLogs"].([]interface{})[0].(map[string]interface{})["scopeLogs"].([]interface{})[0].(map[string]interface{})["logRecords"].([]interface{})[0].(map[string]interface{})["body"] = attributes.AsRaw()

	logRecords := data["resourceLogs"].([]interface{})[0].(map[string]interface{})["scopeLogs"].([]interface{})[0].(map[string]interface{})["logRecords"].([]interface{})
	if len(logRecords) > 0 {
		t1 := time.Now()
		logRecords[0].(map[string]interface{})["timeUnixNano"] = t1.UnixNano()
		attributes := logRecords[0].(map[string]interface{})["attributes"].([]interface{})
		if len(attributes) > 0 {
			// Iterate through replacementAttributes and append each map individually
			for _, attr := range replacementAttributes {
				if attr.(map[string]interface{})["key"] == "service.name" {
					attr.(map[string]interface{})["value"] = map[string]interface{}{
						"stringValue": "testBiqHost",
					}
				} else if attr.(map[string]interface{})["key"] == "telemetry.sdk.name" {
					attr.(map[string]interface{})["value"] = map[string]interface{}{
						"stringValue": "biqconnector",
					}
				}
				attributes = append(attributes, attr)
			}
			//isEventAttr := map[string]interface{}{
			//	"key": "appd.isevent",
			//	"value": map[string]interface{}{
			//		"boolValue": true,
			//	},
			//}
			//attributes = append(attributes, isEventAttr)
			// add span attributes
			for key, value := range spanAttr.AsRaw() {
				attr := map[string]interface{}{
					"key": key,
					"value": map[string]interface{}{
						"stringValue": value,
					},
				}
				attributes = append(attributes, attr)
			}
			logRecords[0].(map[string]interface{})["attributes"] = attributes
		}

		logger.Info("replacementAttributes: ", zap.Any("attributes", attributes))
	}
	// Marshal the modified map back into JSON
	modifiedStr, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error marshaling modified data to JSON:", err)
		return nil, err
	}
	//modifiedJSONStr := strings.Replace(jsonStr, "{\"stringValue\":\"hello world\"}", `"`+test+`"`, 1)
	// Replace the "body" field in the JSON string with the marshaled myStruct JSON
	//modifiedJSONStr := strings.Replace(jsonStr, `"body":null`, `"body":`+string(bodyJSON), 1)

	logger.Info("modifiedJSONStr:", zap.Any("event", modifiedStr))
	decoder := &plog.JSONUnmarshaler{}
	event, err := decoder.UnmarshalLogs(modifiedStr)
	if err != nil {
		fmt.Println("Error unmarshalling Base64 string:", err)
		return nil, err
	}
	logger.Info("Event generated:", zap.Any("event", event))
	return &event, nil
}

func generateId(logger *zap.Logger, resourceAttr pcommon.Map) string {
	service, _ := resourceAttr.Get("service.name")
	logger.Info("BT service: ", zap.Any("service", service))
	namespace, _ := resourceAttr.Get("service.namespace")
	logger.Info("BT namespace: ", zap.Any("namespace", namespace))
	name, _ := resourceAttr.Get("bt.name")
	logger.Info("BT name: ", zap.Any("bt name", name))
	properties := make(map[string]string)
	properties["bt.name"] = name.AsString()
	properties["service.name"] = service.AsString()
	properties["service.namespace"] = namespace.AsString()

	idContent := make([]string, 0, len(properties))

	for k, v := range properties {
		idContent = append(idContent, fmt.Sprintf("%s:%s", k, v))
	}

	sort.Strings(idContent)
	logger.Info("Sorted idContent: ", zap.Any("idContent", idContent))

	var id uuid.UUID
	if idContent == nil || len(idContent) == 0 {
		id = uuid.NewV1()
	} else {
		var longString strings.Builder
		first := true
		for _, content := range idContent {
			if !first {
				longString.WriteString(":")
			}
			longString.WriteString(content)
			first = false
		}
		logger.Info("idContent to longString: ", zap.Any("longString", longString.String()))
		//var err error
		//bytes := []byte(longString.String())
		idPricompact := uuid.NewV5(uuid.NamespaceDNS, longString.String())
		//if err != nil {
		//	logger.Info("Error while generating bt id: ", zap.Any("err", err))
		//}
		//id, _ = uuid.FromBytes(bytes)
		// Parse the UUID string
		id, err := uuid.FromString(idPricompact.String())
		if err != nil {
			logger.Info("Error while parsing uuid: ", zap.Any("err", err))
			return id.String()
		}

		// Extract the 8-byte MSB and LSB from the UUID
		msb := id[:8]
		lsb := id[8:]

		// Create a byte buffer and add MSB and LSB to it
		buffer := make([]byte, 0, 16)
		buffer = append(buffer, msb...)
		buffer = append(buffer, lsb...)

		// Generate a 64-bit encoded string from the byte buffer
		encodedId := base64.StdEncoding.EncodeToString(buffer)
		compactedUuid := strings.Split(encodedId, "=")
		logger.Info("BT id is generated successfully: ", zap.Any("bt-id", compactedUuid[0]))
		return compactedUuid[0]
	}

	logger.Info("BT id generateId: ", zap.Any("bt-id", id))
	return id.String()
}

func ToJSON(obj interface{}) string {
	res, err := json.Marshal(obj)
	if err != nil {
		panic("error with json serialization " + err.Error())
	}
	return string(res)
}

// Helper function to replace a field in a JSON string
func replaceJSONField(jsonStr, oldField, newField string) string {
	return jsonStr[:strings.Index(jsonStr, oldField)] + newField + jsonStr[strings.Index(jsonStr, oldField)+len(oldField):]
}

func getTraceID(td ptrace.Traces) (pcommon.TraceID, error) {
	rss := td.ResourceSpans()
	if rss.Len() == 0 {
		return pcommon.NewTraceIDEmpty(), errors.New("no resource spans are present")
	}

	ilss := rss.At(0).ScopeSpans()
	if ilss.Len() == 0 {
		return pcommon.NewTraceIDEmpty(), errors.New("no scope spans are present")
	}

	spans := ilss.At(0).Spans()
	if spans.Len() == 0 {
		return pcommon.NewTraceIDEmpty(), errors.New("no trace id is present")
	}

	return spans.At(0).TraceID(), nil
}
