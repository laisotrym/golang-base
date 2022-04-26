package safeweb_lib_utils

import (
    "encoding/json"
    "sort"
)

func MapToQueryParam(m map[string]string) string {
    keys := make([]string, len(m))
    idx := 0
    for key := range m {
        keys[idx] = key
        idx++
    }
    sort.Strings(keys)
    data := keys[0] + "=" + m[keys[0]]
    for i := 1; i < len(keys); i++ {
        data += "&" + keys[i] + "=" + m[keys[i]]
    }
    return data
}

func ConvertToMap(in interface{}) (map[string]interface{}, error) {
    var out map[string]interface{}
    inJson, _ := json.Marshal(in)
    err := json.Unmarshal(inJson, &out)
    
    return out, err
}
