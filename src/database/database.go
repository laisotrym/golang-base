package database

import "fmt"

// MySQLConfig schema
type MySQLConfig struct {
	Host     string `json:"host" mapstructure:"host" yaml:"host"`
	Database string `json:"database" mapstructure:"database" yaml:"database"`
	Port     int    `json:"port" mapstructure:"port" yaml:"port"`
	Username string `json:"username" mapstructure:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" yaml:"password"`
	Options  string `json:"options" mapstructure:"options" yaml:"options"`
}

// String return MySQL connection url
func (m MySQLConfig) String() string {
	return fmt.Sprintf("mysql://%s", m.DSN())
}

// DSN return Data Source Name
func (m MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", m.Username, m.Password, m.Host, m.Port, m.Database, m.Options)
}

// PostgreSQLConfig hold config for postgresql
type PostgreSQLConfig struct {
	Host     string `json:"host" mapstructure:"host" yaml:"host"`
	Database string `json:"database" mapstructure:"database" yaml:"database"`
	Port     int    `json:"port" mapstructure:"port" yaml:"port"`
	Username string `json:"username" mapstructure:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" yaml:"password"`
	Options  string `json:"options" mapstrucuture:"options" yaml:"options"`
}

// String return Postgres connection string
// TODO: support ssl
func (m PostgreSQLConfig) String() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s%s", m.Username, m.Password, m.Host, m.Port, m.Database, m.Options)
}

// MySQLDefaultConfig return default config for mysql, usually use on development
func MySQLDefaultConfig() MySQLConfig {
	return MySQLConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "sample",
		Username: "root",
		Password: "khongbiet",
		Options:  "",
	}
}
