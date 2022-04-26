package config

import (
    "fmt"
    
    "github.com/go-redis/redis/v8"
    "github.com/spf13/viper"
)

type RedisConfig struct {
    Addr     string `json:"addr" yaml:"addr"`
    Password string `json:"password" yaml:"password"`
    Db       int    `json:"db" yaml:"db"`
}

func DefaultRedisConfig() RedisConfig {
    return RedisConfig{
        Addr:     "localhost:6379",
        Password: "",
        Db:       0,
    }
}

func InitRedisConfig(vip *viper.Viper) (cnf RedisConfig) {
    cnf = DefaultRedisConfig()
    if vip != nil {
        if err := vip.Unmarshal(&cnf); err != nil {
            fmt.Printf("unable to decode into config struct, %v", err)
        }
    }
    return
}

func NewRedisClient(cnf *RedisConfig) (ins *redis.Client) {
    ins = redis.NewClient(&redis.Options{
        Addr:     cnf.Addr,
        Password: cnf.Password,
        DB:       cnf.Db,
    })
    
    return
}
