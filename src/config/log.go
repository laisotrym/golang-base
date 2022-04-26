package config

import (
    "fmt"
    
    "github.com/spf13/viper"
    "safeweb.app/log"
)

func InitLogConfig(vip *viper.Viper) log.Config {
    cnf := log.DefaultConfig()
    if vip != nil {
        if err := vip.Unmarshal(&cnf); err != nil {
            fmt.Printf("unable to decode into config struct, %v", err)
        }
    }
    return cnf
}
