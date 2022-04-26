package config

import (
	"github.com/spf13/viper"
	"safeweb.app/database"
	safeweb_lib_helper "safeweb.app/service/safeweb_lib/helper"
)

type AppConfig struct {
	// Log              log.Config                        `env:"PromoLog" json:"log" yaml:"log"`
	// Server           server.Config                     `env:"PromoServer" json:"server" yaml:"server"`
	// Client           map[ClientAlias]ClientServiceAddr `env:"ClientService" json:"client" yaml:"client"`
	MySQL            database.MySQLConfig `env:"PromoMySQL" json:"mysql" yaml:"mysql"`
	Redis            RedisConfig          `env:"PromoRedis" json:"redis" yaml:"redis"`
	MigrationsFolder string               `env:"PromoMigrationsFolder" default:"file://migrations" json:"migrations_folder" yaml:"migrations_folder"`
}

func InitAppConfig(vip *viper.Viper, df *viper.Viper) *AppConfig {
	viperMySql := vip.Sub("mysql")
	if viperMySql == nil {
		viperMySql = df.Sub("mysql")
	}
	viperRedis := vip.Sub("redis")
	if viperRedis == nil {
		viperRedis = df.Sub("redis")
	}
	return &AppConfig{
		// Log:              InitLogConfig(vip.Sub("log")),
		// Server:           InitServerConfig(vip.Sub("server")),
		// Client:           InitClientConfig(vip.Sub("client")),
		MySQL:            InitMySQLConfig(viperMySql),
		Redis:            InitRedisConfig(viperRedis),
		MigrationsFolder: safeweb_lib_helper.DefaultString(safeweb_lib_helper.DefaultString(vip.GetString("migrations_folder"), df.GetString("migrations_folder")), "file://migrations"),
	}
}

type AppAlias string

const (
	AppAliasDefault     AppAlias = "service_default"
	AppAliasBackend     AppAlias = "service_backend"
	AppAliasRule        AppAlias = "service_rule"
	AppAliasTransaction AppAlias = "service_transaction"
	AppAliasEngine      AppAlias = "service_engine"
)

var AppNameDict = map[AppAlias]string{
	AppAliasDefault:     "Core Engine",
	AppAliasBackend:     "Core Backend",
	AppAliasRule:        "Core Rule",
	AppAliasTransaction: "Core Transaction",
	AppAliasEngine:      "Core Engine",
}

func (e AppAlias) ToString() string {
	return string(e)
}

func (e AppAlias) Title(s string) string {
	switch s {
	case string(AppAliasDefault):
		return AppNameDict[AppAliasDefault]
	case string(AppAliasBackend):
		return AppNameDict[AppAliasBackend]
	case string(AppAliasRule):
		return AppNameDict[AppAliasRule]
	case string(AppAliasTransaction):
		return AppNameDict[AppAliasTransaction]
	case string(AppAliasEngine):
		return AppNameDict[AppAliasEngine]
	default:
		return ""
	}
}
