package safeweb_lib_helper

import (
    "crypto"
    "fmt"
    "reflect"
)

// GetCacheKey returns the cache key for the given key object by computing a
// checksum of key struct
func HashIfObj(object interface{}) string {
    switch object.(type) {
    case string:
        return object.(string)
    default:
        return Checksum(object)
    }
}

// checksum hashes a given object into a string
func Checksum(object interface{}) string {
    digested := crypto.MD5.New()
    _, _ = fmt.Fprint(digested, reflect.TypeOf(object))
    _, _ = fmt.Fprint(digested, object)
    hash := digested.Sum(nil)
    return fmt.Sprintf("%x", hash)
}
