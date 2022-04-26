package safeweb_lib_cast

import (
    "fmt"
    "reflect"
    "strconv"
    
    "google.golang.org/grpc/codes"
    
    promo_error "safeweb.app/service/safeweb_lib/error"
)

var int64Type = reflect.TypeOf(int64(0))

func GetInt64(unk interface{}) (int64, error) {
    switch i := unk.(type) {
    case int64:
        return i, nil
    case float32:
        return int64(i), nil
    case float64:
        return int64(i), nil
    case int32:
        return int64(i), nil
    case int:
        return int64(i), nil
    case uint64:
        return int64(i), nil
    case uint32:
        return int64(i), nil
    case uint:
        return int64(i), nil
    case string:
        return strconv.ParseInt(i, 10, 64)
    default:
        v := reflect.ValueOf(unk)
        v = reflect.Indirect(v)
        if v.Type().ConvertibleTo(int64Type) {
            fv := v.Convert(int64Type)
            return fv.Int(), nil
        } else if v.Type().ConvertibleTo(stringType) {
            sv := v.Convert(stringType)
            s := sv.String()
            return strconv.ParseInt(s, 10, 64)
        } else {
            return 0, promo_error.NewError(int32(codes.Aborted), "500", fmt.Sprintf(errConvertType, v.Type(), "int64"), false)
        }
    }
}

func GetInt64ReflectOnly(unk interface{}) (int64, error) {
    v := reflect.ValueOf(unk)
    v = reflect.Indirect(v)
    if !v.Type().ConvertibleTo(int64Type) {
        return 0, promo_error.NewError(int32(codes.Aborted), "500", fmt.Sprintf(errConvertType, v.Type(), "int64"), false)
    }
    fv := v.Convert(int64Type)
    return fv.Int(), nil
}

func GetInt64SwitchOnly(unk interface{}) (int64, error) {
    switch i := unk.(type) {
    case int64:
        return i, nil
    case float32:
        return int64(i), nil
    case float64:
        return int64(i), nil
    case int32:
        return int64(i), nil
    case int:
        return int64(i), nil
    case uint64:
        return int64(i), nil
    case uint32:
        return int64(i), nil
    case uint:
        return int64(i), nil
    default:
        return 0, promo_error.NewError(int32(codes.Aborted), "500", fmt.Sprintf(errUnexpectedType, "int64"), false)
    }
}
