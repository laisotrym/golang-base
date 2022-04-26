package safeweb_lib_error

import (
    "fmt"
    "os"
    "reflect"

    "google.golang.org/grpc/codes"
    "gopkg.in/yaml.v2"

    "safeweb.app/service/safeweb_lib/helper"
)

// Load from status.yml if the has not initialized yet.
func Load(input interface{}) error {
    statusPath := os.Getenv("safeweb_lib_ERROR_PATH")
    if statusPath != "" {
        if !safeweb_lib_helper.FileExists(statusPath) {
            panic(fmt.Sprintf("%s path invalid", statusPath))
        }
        return Init(statusPath, input)
    } else {
        panic("Environment PROMO_LIB_ERROR_PATH not set")
    }
}

// Init Load statuses from the given config file.
// Init panics if cannot access or error while parsing the config file.
func Init(filePath string, data interface{}) error {
    f, err := os.Open(filePath)
    if err != nil {
        panic(err.Error())
    }
    defer func() {
        _ = f.Close()
    }()
    if err := yaml.NewDecoder(f).Decode(data); err != nil {
        return err
    }

    return checkNotNilData(data)
}

func checkNotNilData(v interface{}) error {
    strut := reflect.Indirect(reflect.ValueOf(v)).Interface()
    value := reflect.ValueOf(strut).Elem()
    eType := reflect.Indirect(value)
    if value.Kind() == reflect.Struct {
        for i := 0; i < value.NumField(); i++ {
            child := value.Field(i)
            childType := eType.Type().Field(i)
            catName := childType.Name
            if child.Kind() == reflect.Struct {
                for j := 0; j < child.NumField(); j++ {
                    errName := childType.Type.Field(j).Name
                    val := child.Field(j)
                    isNil := false
                    k := val.Kind()
                    switch k {
                    case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
                        isNil = val.IsNil()
                    default:
                        isNil = val.IsZero()
                    }
                    if isNil == true {
                        promoErr := NewError(int32(codes.Internal), "500", fmt.Sprintf("%s.%s", catName, errName), false)
                        val.Set(reflect.ValueOf(promoErr))
                    }
                }
            }
        }
    }

    return nil
}

func FindGeneralWithCode(code string) *Error {
    statuses, err := GetGeneralStatus()
    if err != nil {
        return nil
    }
    v := reflect.ValueOf(statuses)
    for i := 0; i < v.NumField(); i++ {
        s, ok := v.Field(i).Interface().(*Error)
        if !ok || s == nil {
            continue
        }
        if s.Code == code {
            return s
        }
    }
    return nil
}
