package safeweb_lib_validate

import (
    "reflect"
    "strings"

    "github.com/go-playground/validator"

    promo_error "safeweb.app/service/safeweb_lib/error"
    "safeweb.app/service/safeweb_lib/utils"
)

type Validator interface {
    Struct(v interface{}) (bool, error)
    Variable(in interface{}, rule string) (bool, error)
}

var (
    instance Validator
)

func getInstance() Validator {
    if instance == nil {
        instance = NewValidator()
    }
    return instance
}

func Struct(v interface{}) (bool, error) {
    result, err := getInstance().Struct(v)
    if result {
        return result, nil
    }
    data := err.(validator.ValidationErrors)
    e := data[0]
    tag := e.Tag()
    transData := translate(e)
    fieldName := e.StructField()
    var message string
    if fieldTag := findFieldNameTag(v, fieldName); fieldTag != "" {
        message = strings.Replace(transData, fieldName, fieldTag, 1)
    } else {
        message = strings.Replace(transData, fieldName, safeweb_lib_utils.StringToSnakeCase(fieldName), 1)
    }
    
    state, _ := promo_error.GetGeneralStatus()
    
    switch tag {
    case "required":
        return false, state.MissingRequiredParam.ReplaceMessage(message)
    default:
        return false, state.DataInvalid.ReplaceMessage(message)
    }
}

func findFieldNameTag(v interface{}, fieldName string) string {
    val := reflect.Indirect(reflect.ValueOf(v))
    field, ok := val.Type().FieldByName(fieldName)
    if !ok {
        return ""
    }
    a := field.Tag.Get("field_name")
    return a
}

func Variable(v interface{}, rule string) (bool, error) {
    return getInstance().Variable(v, rule)
}
