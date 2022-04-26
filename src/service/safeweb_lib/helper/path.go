package safeweb_lib_helper

import (
    "os"
    "path"
    "path/filepath"
    "runtime"
)

func CurrentDir() string {
    _, b, _, _ := runtime.Caller(0)
    d := path.Join(path.Dir(b))
    return filepath.Dir(d)
}

func CurrDir() string {
    if dir, err := os.Getwd(); err != nil {
        return ""
    } else {
        return dir
    }
}

func FullPath(fileName string) string {
    return path.Join(CurrDir(), fileName)
}

func FileExists(f string) bool {
    info, err := os.Stat(f)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}
