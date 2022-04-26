package config

import (
    "github.com/go-redis/redis/v8"
    "github.com/gocraft/dbr/v2"
)

func InitConn(conf *AppConfig, initRedis bool, initConn bool) (db *dbr.Connection, res *redis.Client, err error) {
    // Init db connection
    if initConn {
        db, err = DbConnect(conf.MySQL.DSN())
        if err != nil {
            return
        }
    }

    if initRedis {
        res = NewRedisClient(&conf.Redis)
    }

    return
}
