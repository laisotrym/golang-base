package config

import (
    "fmt"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/alicebob/miniredis"
    "github.com/go-redis/redis/v8"
    "github.com/gocraft/dbr/v2"
    "github.com/gocraft/dbr/v2/dialect"
)

type TestStruct struct {
    Obj         interface{}
    Query       string
    EmptyError  bool
    EmptyOutput bool
}

func InitMockConn(initRedis bool) (*dbr.Connection, sqlmock.Sqlmock, *redis.Client, error) {
    db, mock, err := sqlmock.New()
    if err != nil {
        return nil, nil, nil, fmt.Errorf("Error  %v not expect when mock db", err)
    }
    conn := &dbr.Connection{
        DB:            db,
        Dialect:       dialect.MySQL,
        EventReceiver: &dbr.NullEventReceiver{},
    }
    
    if !initRedis {
        return conn, mock, nil, err
    }
    
    mr, err := miniredis.Run()
    if err != nil {
        return nil, nil, nil, fmt.Errorf("Error  %v not expect when mock redis", err)
    }
    red := redis.NewClient(&redis.Options{
        Addr: mr.Addr(),
    })
    return conn, mock, red, err
    
}
