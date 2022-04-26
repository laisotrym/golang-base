package safeweb_lib_validate

import (
    "fmt"
    "log"
    "reflect"
    "regexp"
    "strconv"
    "sync"
    
    "github.com/go-playground/locales/en"
    ut "github.com/go-playground/universal-translator"
    
    viTran "safeweb.app/service/safeweb_lib/validate/vi"
    
    "github.com/go-playground/validator"
)

type ValidatorImpl struct {
    v *validator.Validate
}

var (
    trans ut.Translator
    uni   *ut.UniversalTranslator
    once  sync.Once
)

// New return instance of validate
func NewValidator() *ValidatorImpl {
    
    valid := &ValidatorImpl{}
    once.Do(func() {
        // get translator
        en := en.New()
        uni = ut.New(en, en)
        
        trans, _ = uni.GetTranslator("en")
        valid.v = validator.New()
        err := viTran.RegisterDefaultTranslations(valid.v, trans)
        if err != nil {
            log.Printf("error when register default translation: %+v\n", err)
        }
    })
    _ = valid.v.RegisterValidation("code", validateCommonCode)
    _ = valid.v.RegisterValidation("api", validateUrlApi)
    _ = valid.v.RegisterValidation("optional_len_lte", OptionalStringLenGte)
    _ = valid.v.RegisterValidation("optional_len_gte", OptionalStringLenLte)
    _ = valid.v.RegisterValidation("optional_len", OptionalStringLen)
    
    return valid
}
func (v *ValidatorImpl) Struct(s interface{}) (bool, error) {
    err := v.v.Struct(s)
    return err == nil, err
}
func (v *ValidatorImpl) Variable(in interface{}, rule string) (bool, error) {
    err := v.v.Var(in, rule)
    return err == nil, err
}

func translate(err validator.FieldError) string {
    return err.Translate(trans)
}

func validateCommonCode(fl validator.FieldLevel) bool {
    matched, _ := regexp.Match("^[a-zA-Z0-9_-]*$", []byte(fl.Field().String()))
    return matched
}

func validateUrlApi(fl validator.FieldLevel) bool {
    matched, _ := regexp.Match(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`, []byte(fl.Field().String()))
    return matched
}

func OptionalStringLenGte(fl validator.FieldLevel) bool {
    field := fl.Field()
    fieldType := field.Kind()
    var processString string
    switch fieldType {
    case reflect.String:
        processString = field.String()
        break
    case reflect.Int:
        processString = fmt.Sprintf("%d", field.Interface())
        break
    case reflect.Int64:
        processString = fmt.Sprintf("%d", field.Interface())
        break
    }
    // if field empty
    if len(processString) == 0 || processString == "0" {
        return true
    }
    length, err := strconv.Atoi(fl.Param())
    if err != nil {
        log.Fatal("error validate optional string gte", err)
        return false
    }
    if len([]rune(processString)) > length {
        return false
    }
    return true
}

func OptionalStringLenLte(fl validator.FieldLevel) bool {
    field := fl.Field()
    fieldType := field.Kind()
    var processString string
    switch fieldType {
    case reflect.String:
        processString = field.String()
        break
    case reflect.Int:
        processString = fmt.Sprintf("%d", field.Interface())
        break
    case reflect.Int64:
        processString = fmt.Sprintf("%d", field.Interface())
        break
    }
    // if field empty
    if len(processString) == 0 || processString == "0" {
        return true
    }
    length, err := strconv.Atoi(fl.Param())
    if err != nil {
        log.Fatal("error validate optional string lte", err)
        return false
    }
    if len([]rune(processString)) < length {
        return false
    }
    return true
}

func OptionalStringLen(fl validator.FieldLevel) bool {
    field := fl.Field()
    fieldType := field.Kind()
    var processString string
    switch fieldType {
    case reflect.String:
        processString = field.String()
        break
    case reflect.Int:
        processString = fmt.Sprintf("%d", field.Interface())
        break
    case reflect.Int64:
        processString = fmt.Sprintf("%d", field.Interface())
        break
    }
    // if field empty
    if len(processString) == 0 || processString == "0" {
        return true
    }
    length, err := strconv.Atoi(fl.Param())
    if err != nil {
        log.Fatal("error validate optional string lte", err)
        return false
    }
    if len([]rune(processString)) != length {
        return false
    }
    return true
}
