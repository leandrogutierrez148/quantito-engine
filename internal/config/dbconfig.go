package config

import (
	"github.com/go-sql-driver/mysql"
)

// DBConfig represents Database configuration model
type DBConfig struct {
	Host            string `toml:"host"`
	Username        string `toml:"username"`
	Password        string `toml:"password"`
	Schema          string `toml:"schema"`
	MaxOpenConn     int    `toml:"maxOpenConn"`
	MaxIdleConn     int    `toml:"maxIdleConn"`
	MaxTTLConn      int    `toml:"maxTTLConn"`
	MaxIdleTimeConn int    `toml:"maxIdleTimeConn"`
}

// ToMysqlConfig converts DBConfig to mysql.Config
func (dbConfig *DBConfig) ToMysqlConfig() *mysql.Config {
	config := mysql.NewConfig()
	config.User = dbConfig.Username
	config.Passwd = dbConfig.Password
	config.Net = "tcp"
	config.Addr = dbConfig.Host
	config.DBName = dbConfig.Schema

	return config
}
