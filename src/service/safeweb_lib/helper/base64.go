package safeweb_lib_helper

import (
    "encoding/base64"
)

func EncodeBase64(str string) string {
    return base64.StdEncoding.EncodeToString([]byte(str))
}

func DecodeBase64(str string) ([]byte, error) {
    return base64.StdEncoding.DecodeString(str)
}
