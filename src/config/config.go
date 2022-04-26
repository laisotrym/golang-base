package config

import (
	"fmt"

	"github.com/spf13/viper"
	"safeweb.app/database"
	"safeweb.app/log"

	"safeweb.app/server"
	safeweb_lib_helper "safeweb.app/service/safeweb_lib/helper"
)

type Config struct {
	TimeZone         string                            `json:"time_zone" yaml:"time_zone" env:"SAFEWEB_TIMEZONE" default:"UTC"`
	Log              log.Config                        `json:"log" yaml:"log"`
	Server           server.Config                     `json:"server" yaml:"server"`
	MySQL            database.MySQLConfig              `json:"mysql" yaml:"mysql"`
	Redis            RedisConfig                       `json:"redis" yaml:"redis"`
	MigrationsFolder string                            `json:"migrations_folder" yaml:"migrations_folder" env:"SAFEWEB_MIGRATIONS_FOLDER" default:"file://migrations"`
	Client           map[ClientAlias]ClientServiceAddr `json:"client" yaml:"client"`

	App map[AppAlias]*AppConfig `json:"app" yaml:"app"`

	PrivateKey map[string]string `json:"private_key" yaml:"privateKey"`
	SuperKey   map[string]string `json:"super_key" yaml:"superKey"`
	Webroot    string            `yaml:"webroot"`
}

func Load() (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./bin")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Read config from file %s error: %s\n", viper.ConfigFileUsed(), err.Error())
		return nil, err
	} else {
		fmt.Printf("Read config from file %s\n", viper.ConfigFileUsed())
	}

	cfg := &Config{
		TimeZone:         safeweb_lib_helper.DefaultString(viper.GetString("time_zone"), "UTC"),
		Log:              InitLogConfig(viper.Sub("log")),
		MySQL:            InitMySQLConfig(viper.Sub("mysql")),
		Server:           InitServerConfig(viper.Sub("server")),
		Redis:            InitRedisConfig(viper.Sub("redis")),
		MigrationsFolder: safeweb_lib_helper.DefaultString(viper.GetString("migrations_folder"), "file://migrations"),
		Client:           InitClientConfig(viper.Sub("client")),

		App: make(map[AppAlias]*AppConfig),

		PrivateKey: InitApiKeyConfig(viper.Sub("privateKey")),

		SuperKey: InitApiKeyConfig(viper.Sub("superKey")),
		Webroot:  viper.GetString("webroot"),
	}

	viperRoot := viper.GetViper()
	cfg.App[AppAliasDefault] = InitAppConfig(viperRoot, viperRoot)
	viperBackend := viper.Sub(AppAliasBackend.ToString())
	if viperBackend != nil {
		cfg.App[AppAliasBackend] = InitAppConfig(viperBackend, viperRoot)
	}
	viperRule := viper.Sub(AppAliasRule.ToString())
	if viperRule != nil {
		cfg.App[AppAliasRule] = InitAppConfig(viperRule, viperRoot)
	}
	viperTrans := viper.Sub(AppAliasTransaction.ToString())
	if viperTrans != nil {
		cfg.App[AppAliasTransaction] = InitAppConfig(viperTrans, viperRoot)
	}
	viperEngine := viper.Sub(AppAliasEngine.ToString())
	if viperEngine != nil {
		cfg.App[AppAliasEngine] = InitAppConfig(viperEngine, viperRoot)
	}

	return cfg, nil
}
