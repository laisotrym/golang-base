package config

import (
    "fmt"

    "github.com/spf13/viper"
)

func InitApiKeyConfig(vip *viper.Viper) map[string]string {
    cnf := make(map[string]string)
    if vip != nil {
        if err := vip.Unmarshal(&cnf); err != nil {
            fmt.Printf("unable to decode into config struct, %v", err)
        }
    }
    return cnf
}
