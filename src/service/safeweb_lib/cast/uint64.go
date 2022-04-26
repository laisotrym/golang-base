package safeweb_lib_cast

import (
    "fmt"
    "reflect"
    "strconv"
    
    "google.golang.org/grpc/codes"
    
    promo_error "safeweb.app/service/safeweb_lib/error"
)

var uint64Type = reflect.TypeOf(uint64(0))

func GetUInt64(unk interface{}) (uint64, error) {
    switch i := unk.(type) {
    case uint64:
        return i, nil
    case float32:
        return uint64(i), nil
    case float64:
        return uint64(i), nil
    case int32:
        return uint64(i), nil
    case int:
        return uint64(i), nil
    case int64:
        return uint64(i), nil
    case uint32:
        return uint64(i), nil
    case uint:
        return uint64(i), nil
    case string:
        return strconv.ParseUint(i, 10, 64)
    default:
        v := reflect.ValueOf(unk)
        v = reflect.Indirect(v)
        if v.Type().ConvertibleTo(uint64Type) {
            fv := v.Convert(uint64Type)
            return fv.Uint(), nil
        } else if v.Type().ConvertibleTo(stringType) {
            sv := v.Convert(stringType)
            s := sv.String()
            return strconv.ParseUint(s, 10, 64)
        } else {
            return 0, promo_error.NewError(int32(codes.Aborted), "500", fmt.Sprintf(errConvertType, v.Type(), "uint64"), false)
        }
    }
}

func GetUInt64ReflectOnly(unk interface{}) (uint64, error) {
    v := reflect.ValueOf(unk)
    v = reflect.Indirect(v)
    if !v.Type().ConvertibleTo(uint64Type) {
        return 0, promo_error.NewError(int32(codes.Aborted), "500", fmt.Sprintf(errConvertType, v.Type(), "uint64"), false)
    }
    fv := v.Convert(uint64Type)
    return fv.Uint(), nil
}

func GetUInt64SwitchOnly(unk interface{}) (uint64, error) {
    switch i := unk.(type) {
    case uint64:
        return i, nil
    case float32:
        return uint64(i), nil
    case float64:
        return uint64(i), nil
    case int32:
        return uint64(i), nil
    case int:
        return uint64(i), nil
    case int64:
        return uint64(i), nil
    case uint32:
        return uint64(i), nil
    case uint:
        return uint64(i), nil
    default:
        return 0, promo_error.NewError(int32(codes.Aborted), "500", fmt.Sprintf(errUnexpectedType, "uint64"), false)
    }
}
