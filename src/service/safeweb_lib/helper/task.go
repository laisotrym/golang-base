package safeweb_lib_helper

import (
    "errors"
    "time"
    
    "github.com/spf13/viper"
)

type AsyncFunction func() (interface{}, error)
type AsyncResult func(timeout ...time.Duration) (interface{}, error)

func AsyncTask(f AsyncFunction) AsyncResult {
    var result interface{}
    var err error
    
    c := make(chan struct{}, 1)
    
    go func() {
        defer close(c)
        result, err = f()
    }()
    
    return func(timeout ...time.Duration) (interface{}, error) {
        if len(timeout) == 1 {
            select {
            case <-c:
                return result, err
            case <-time.After(timeout[0]):
                return nil, errors.New("async timeout")
            }
        } else {
            <-c
            return result, err
        }
    }
}

type VoidFunc = func()

func DeferFunc(f VoidFunc) (result bool, err error) {
    defer func() {
        if r := recover(); r != nil {
            switch x := r.(type) {
            case string:
                err = errors.New(x)
            case error:
                err = x
            default:
                err = errors.New("unknown panic")
            }
        }
    }()
    f()
    return true, nil
}

func HandleAsyncErrors(asyncResults []AsyncResult) {
    for _, asyncResult := range asyncResults {
        _, e := asyncResult()
        if e != nil {
            panic(e)
        }
    }
}

func ShouldForwardOrder() bool {
    return viper.GetBool("FORWARD_ORDER") == true
}
