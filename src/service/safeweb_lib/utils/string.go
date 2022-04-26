package safeweb_lib_utils

import (
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func StringToSnakeCase(str string) string {
    snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
    snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
    return strings.ToLower(snake)
}

func Trim(obj interface{}) interface{} {
    original := reflect.ValueOf(obj)
    copy := reflect.New(original.Type()).Elem()
    trimRecursive(copy, original)
    return copy.Interface()
}

func trimRecursive(copy, original reflect.Value) {
    switch original.Kind() {
    case reflect.Ptr:
        originalValue := original.Elem()
        if !originalValue.IsValid() {
            return
        }
        copy.Set(reflect.New(originalValue.Type()))
        trimRecursive(copy.Elem(), originalValue)
    case reflect.Interface:
        originalValue := original.Elem()
        copyValue := reflect.New(originalValue.Type()).Elem()
        trimRecursive(copyValue, originalValue)
        copy.Set(copyValue)
    case reflect.Struct:
        for i := 0; i < original.NumField(); i++ {
            trimRecursive(copy.Field(i), original.Field(i))
        }
    case reflect.Slice, reflect.Array:
        if original.Len() > 0 && original.Cap() > 0 {
            copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
            for i := 0; i < original.Len(); i++ {
                trimRecursive(copy.Index(i), original.Index(i))
            }
        }
    case reflect.Map:
        copy.Set(reflect.MakeMap(original.Type()))
        for _, key := range original.MapKeys() {
            originalValue := original.MapIndex(key)
            copyValue := reflect.New(originalValue.Type()).Elem()
            trimRecursive(copyValue, originalValue)
            copy.SetMapIndex(key, copyValue)
        }
    case reflect.String:
        trimString := strings.TrimSpace(original.Interface().(string))
        copy.SetString(trimString)
    default:
        copy.Set(original)
    }
}

func SliceInt64ToString(values []int64) string {
    var valuesText []string
    
    for i := range values {
        number := values[i]
        text := strconv.FormatInt(number, 10)
        valuesText = append(valuesText, text)
    }
    
    return strings.Join(valuesText, ",")
}

func SliceIntToString(values []int) string {
    var valuesText []string
    
    for i := range values {
        number := values[i]
        text := strconv.Itoa(number)
        valuesText = append(valuesText, text)
    }
    
    return strings.Join(valuesText, ",")
}

func SplitTrimSpace(s, sep string) (r []string) {
    r = strings.Split(s, sep)
    for i := range r {
        r[i] = strings.TrimSpace(r[i])
    }
    
    return
}

func DeleteEmpty(s []string) (r []string) {
    for _, i := range s {
        if strings.TrimSpace(i) != "" {
            r = append(r, strings.TrimSpace(i))
        }
    }
    return
}

func RemoveEmpty(slice *[]string) {
    i := 0
    p := *slice
    for _, entry := range p {
        if strings.Trim(entry, " ") != "" {
            p[i] = entry
            i++
        }
    }
    *slice = p[0:i]
}
