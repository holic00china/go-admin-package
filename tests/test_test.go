package tests

import (
	"testing"
	"time"

	"go-admin-package/modules/config"
	"go-admin-package/tests/frameworks/gin"
)

func TestBlackBoxTestSuitOfBuiltInTables(t *testing.T) {
	BlackBoxTestSuitOfBuiltInTables(t, gin.NewHandler, config.DatabaseList{
		"default": {
			Host:            "127.0.0.1",
			Port:            "3306",
			User:            "root",
			Pwd:             "root",
			Name:            "go-admin-package-test",
			MaxIdleConns:    50,
			MaxOpenConns:    150,
			ConnMaxLifetime: time.Hour,
			ConnMaxIdleTime: 0,
			Driver:          config.DriverMysql,
		},
	})
}
