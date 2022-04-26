package safeweb_lib_helper

import (
    "strings"

    "safeweb.app/service/safeweb_lib/types/set"
)

func IsEmpty(str string) bool {
    trim := strings.TrimSpace(str)
    isEmpty := len(trim) == 0
    return isEmpty
}

func IsIntersection(a []string, b []string) bool {
    setA := safeweb_lib_types_set.NewStringSet(a...)
    setB := safeweb_lib_types_set.NewStringSet(b...)
    intersection := safeweb_lib_types_set.Intersection(setA, setB)
    return intersection != nil && intersection.Size() > 0
}
