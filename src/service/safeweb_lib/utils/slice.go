package safeweb_lib_utils

import (
    "fmt"
    "reflect"
    "strconv"

    "github.com/pkg/errors"

    "safeweb.app/model/enum"
)

func FindFirstDuplicateString(arr []string) (string, bool) {
    check := make(map[string]int)
    
    for _, s := range arr {
        if check[s] > 0 {
            return s, true
        }
        check[s] = 1
    }
    return "", false
}

func FindMissingString(arr []string, keys ...string) []string {
    miss := make([]string, 0)
    for _, a := range arr {
        found := false
        for _, k := range keys {
            if a == k {
                found = true
                break
            }
        }
        if !found {
            miss = append(miss, a)
        }
    }
    return miss
}

func FindOutsideKeys(arr []string, keys ...string) []string {
    outside := make([]string, 0)
    for _, k := range keys {
        found := false
        for _, a := range arr {
            if a == k {
                found = true
            }
        }
        if !found {
            outside = append(outside, k)
        }
    }
    return outside
}

func RemoveDuplicatedInt64(arr []int64) []int64 {
    m := make(map[int64]bool)
    res := []int64{}
    for _, value := range arr {
        if _, ok := m[value]; !ok {
            m[value] = true
            res = append(res, value)
        }
    }
    return res
}

// Get complement slice between 2 slices, that means get all values exists in arr1 but not arr2
// Example: arr1 = [1, 2, 3], arr2 = [3, 4] => result = [1, 2]
func GetComplementSliceInt64(arr1 []int64, arr2 []int64) []int64 {
    m2 := make(map[int64]bool)
    for _, v2 := range arr2 {
        if _, ok := m2[v2]; !ok {
            m2[v2] = true
        }
    }
    
    res := []int64{}
    for _, v1 := range arr1 {
        if _, ok := m2[v1]; !ok {
            res = append(res, v1)
        }
    }
    
    return res
}

func CheckExists(items interface{}, s string) (int, bool) {
    if reflect.TypeOf(items).Kind() == reflect.Slice {
        objs := reflect.ValueOf(items)
        if objs.Len() == 0 {
            return 0, false
        }
        for i := 0; i < objs.Len(); i++ {
            item := objs.Index(i).Interface()
            
            switch item.(type) {
            case enum.IEnum:
                if item.(enum.IEnum).ToString() == s {
                    return i, true
                }
            case uint64:
                itemValueStr := fmt.Sprintf("%v", item)
                if itemValueStr == s {
                    return i, true
                }
            }
        }
    }
    return -1, false
}

func ConvertToNumberArr(arr []string) ([]int64, error) {
    var numberArr = make([]int64, 0, len(arr))
    for _, value := range arr {
        numValue, err := strconv.ParseInt(value, 10, 64)
        if err != nil {
            return nil, errors.New("data's format is invalid")
        }
        numberArr = append(numberArr, numValue)
    }
    return numberArr, nil
}

func MakeSliceFromInput(input interface{}) ([]interface{}, error) {
    s := reflect.ValueOf(input)
    if s.Kind() != reflect.Slice {
        return nil, errors.New("Invalid input request")
    }
    slices := make([]interface{}, s.Len())
    for i := 0; i < s.Len(); i++ {
        slices[i] = s.Index(i).Interface()
    }
    return slices, nil
}
