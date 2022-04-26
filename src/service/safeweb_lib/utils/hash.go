package safeweb_lib_utils

import (
    "crypto/md5"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

func HashMD5(s string) string {
    hasher := md5.New()
    hasher.Write([]byte(s))
    return hex.EncodeToString(hasher.Sum(nil))
}

func HashSHA256(s string) string {
    hash := sha256.Sum256([]byte(s))
    return fmt.Sprintf("%x", hash[:])
}

func Hash(s string, hashType string) string {
    switch hashType {
    case "SHA256":
        return HashSHA256(s)
    case "MD5":
        return HashMD5(s)
    }
    return ""
}
