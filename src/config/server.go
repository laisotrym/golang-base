package config

import (
    "fmt"
    
    "github.com/spf13/viper"
    
    "safeweb.app/server"
)

func InitServerConfig(vip *viper.Viper) server.Config {
    cnf := server.DefaultConfig()
    if vip != nil {
        if err := vip.Unmarshal(&cnf); err != nil {
            fmt.Printf("unable to decode into config struct, %v", err)
        }
    }
    return cnf
}
