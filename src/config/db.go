package config

import (
	"fmt"

	"github.com/gocraft/dbr/v2"
	"github.com/spf13/viper"
    "safeweb.app/database"
)

func InitMySQLConfig(vip *viper.Viper) database.MySQLConfig {
    cnf := database.MySQLDefaultConfig()
    if vip != nil {
        if err := vip.Unmarshal(&cnf); err != nil {
            fmt.Printf("unable to decode into config struct, %v", err)
        }
    }
    return cnf
}

func DbConnect(dbURL string) (*dbr.Connection, error) {
    if conn, err := dbr.Open("mysql", dbURL, nil); err != nil {
        return nil, err
    } else {
        conn.SetMaxIdleConns(10)
        conn.SetMaxOpenConns(100)
        
        return conn, nil
    }
}
