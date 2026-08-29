package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMysqlConfig(t *testing.T) {
	dbConfig := DBConfig{
		Username: "user",
		Password: "123456",
		Schema:   "mydb",
		Host:     "dbhost:3308",
	}

	config := dbConfig.ToMysqlConfig()

	assert.Equal(t, "user", config.User)
	assert.Equal(t, "123456", config.Passwd)
	assert.Equal(t, "mydb", config.DBName)
	assert.Equal(t, "dbhost:3308", config.Addr)
	assert.Equal(t, "tcp", config.Net)
}
