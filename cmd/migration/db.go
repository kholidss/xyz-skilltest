package migration

import (
	"fmt"
	"github.com/kholidss/xyz-skilltest/pkg/database/mysql"

	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
)

func MigrateDatabase() {
	cfg, err := config.LoadAllConfigs()

	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to load configuration file: %v", err))
	}

	mysql.DatabaseMigration(cfg)
}
