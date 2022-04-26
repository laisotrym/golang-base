package safeweb_lib_helper

import (
    "encoding/json"
    
    "github.com/go-openapi/swag"
)

func IsZero(items ...interface{}) bool {
    for _, item := range items {
        if !swag.IsZero(item) {
            return false
        }
    }
    return true
}

func DefaultString(str string, defaultValue string) string {
    if IsZero(str) {
        return defaultValue
    }
    return str
}

func Default(item interface{}, defaultValue interface{}) interface{} {
    if IsZero(item) {
        return defaultValue
    }
    
    return item
}

func ErrorOrNil(e interface{}) error {
    if e == nil {
        return nil
    }
    return e.(error)
}

func BytesToMap(data []byte) interface{} {
    var raw map[string]interface{}
    if err := json.Unmarshal(data, &raw); err != nil {
        return data
    }
    return raw
}
