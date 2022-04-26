package gatewayopt

import (
	"reflect"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// JSONPbMarshaler is a Marshaler that translates GRPC response to HTTP JSON response that conforms standard.
// Example:
// message ListBooksResponse {
//	 repeated Book books = 1 [(gogoproto.moretags)="response_field:\"data\""];
//   string message = 2;
// }
// will be translated to JSON message:
// {
//	 "data": [
//	   // list of book here
//	   {...}
//	 ],
//   "message": "a message"
// }
type JSONPbMarshaler struct {
	*runtime.JSONPb
	removeOriRespDataFields bool
}

// Marshal marshals "v" into JSON, following standardised format.
func (j JSONPbMarshaler) Marshal(v interface{}) ([]byte, error) {
	protoMsg, ok := v.(proto.Message)
	if !ok {
		return j.JSONPb.Marshal(v)
	}

	type responseData struct {
		Name      string
		JSONTag   string
		OmitEmpty bool
	}

	respDataFields := make(map[string][]*responseData)

	respMap := make(map[string]interface{})

	s := reflect.ValueOf(protoMsg).Elem()
	for i := 0; i < s.NumField(); i++ {
		value := s.Field(i)
		valueField := s.Type().Field(i)
		if strings.HasPrefix(valueField.Name, "XXX_") {
			continue
		}

		// this is not a protobuf field
		if valueField.Tag.Get("protobuf") == "" && valueField.Tag.Get("protobuf_oneof") == "" {
			continue
		}

		jsonTag := valueField.Tag.Get("json")
		var omitEmpty bool
		if strings.HasSuffix(jsonTag, ",omitempty") {
			omitEmpty = true
			jsonTag = strings.TrimSuffix(jsonTag, ",omitempty")
		}

		if _, ok := value.Interface().(proto.Message); ok && value.IsNil() {
			// to avoid err "Marshal called with nil"
			// https://github.com/golang/protobuf/blob/v1.4.3/jsonpb/encode.go#L88
			respMap[jsonTag] = nil
		} else {
			respMap[jsonTag] = value.Interface()
		}

		responseFieldName := valueField.Tag.Get("response_field")
		if responseFieldName != "" {
			// skip if tag name is same as response attribute json tag
			if _, ok := respMap[responseFieldName]; ok {
				continue
			}
			respDataFields[responseFieldName] = append(respDataFields[responseFieldName], &responseData{
				Name:      valueField.Name,
				JSONTag:   jsonTag,
				OmitEmpty: omitEmpty,
			})
		}
	}

	for name, fields := range respDataFields {
		if len(fields) == 1 {
			field := fields[0]
			respMap[name] = respMap[field.JSONTag]
			if j.removeOriRespDataFields {
				delete(respMap, field.JSONTag)
			}
		} else if len(fields) > 1 {
			wrappedData := make(map[string]interface{})
			for _, field := range fields {
				wrappedData[field.JSONTag] = respMap[field.JSONTag]
				if j.removeOriRespDataFields {
					delete(respMap, field.JSONTag)
				}
			}
			respMap[name] = wrappedData
		}
	}

	return j.JSONPb.Marshal(respMap)
}
