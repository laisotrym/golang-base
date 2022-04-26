package safeweb_lib_cast

import (
    "fmt"
    "math"
    "reflect"
    "strconv"
    
    "google.golang.org/grpc/codes"
    
    promo_error "safeweb.app/service/safeweb_lib/error"
)

var float64Type = reflect.TypeOf(float64(0))

func GetFloat64(unk interface{}) (float64, error) {
    switch i := unk.(type) {
    case float64:
        return i, nil
    case float32:
        return float64(i), nil
    case int64:
        return float64(i), nil
    case int32:
        return float64(i), nil
    case int:
        return float64(i), nil
    case uint64:
        return float64(i), nil
    case uint32:
        return float64(i), nil
    case uint:
        return float64(i), nil
    case string:
        return strconv.ParseFloat(i, 64)
    default:
        v := reflect.ValueOf(unk)
        v = reflect.Indirect(v)
        if v.Type().ConvertibleTo(float64Type) {
            fv := v.Convert(float64Type)
            return fv.Float(), nil
        } else if v.Type().ConvertibleTo(stringType) {
            sv := v.Convert(stringType)
            s := sv.String()
            return strconv.ParseFloat(s, 64)
        } else {
            return math.NaN(), promo_error.NewError(int32(codes.Aborted), "500", fmt.Sprintf(errConvertType, v.Type(), "float64"), false)
        }
    }
}

func GetFloat64ReflectOnly(unk interface{}) (float64, error) {
    v := reflect.ValueOf(unk)
    v = reflect.Indirect(v)
    if !v.Type().ConvertibleTo(float64Type) {
        return math.NaN(), promo_error.NewError(int32(codes.Aborted), "500", fmt.Sprintf(errConvertType, v.Type(), "float64"), false)
    }
    fv := v.Convert(float64Type)
    return fv.Float(), nil
}

func GetFloat64SwitchOnly(unk interface{}) (float64, error) {
    switch i := unk.(type) {
    case float64:
        return i, nil
    case float32:
        return float64(i), nil
    case int64:
        return float64(i), nil
    case int32:
        return float64(i), nil
    case int:
        return float64(i), nil
    case uint64:
        return float64(i), nil
    case uint32:
        return float64(i), nil
    case uint:
        return float64(i), nil
    default:
        return math.NaN(), promo_error.NewError(int32(codes.Aborted), "500", fmt.Sprintf(errUnexpectedType, "float64"), false)
    }
}
