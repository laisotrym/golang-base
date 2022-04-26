package safeweb_lib_helper

import (
    "reflect"
)

// check if item exist in array
func ItemExists(list interface{}, item interface{}) bool {
    arr := reflect.ValueOf(list)
    for i := 0; i < arr.Len(); i++ {
        if arr.Index(i).Interface() == item {
            return true
        }
    }
    
    return false
}
