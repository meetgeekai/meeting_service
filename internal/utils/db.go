package utils

import (
	"fmt"

	"github.com/meetgeekai/go-common/utils"
)

func PrepareMySQLDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC",
		utils.GetEnv[string]("MYSQL_USER"),
		utils.GetEnv[string]("MYSQL_PASS"),
		utils.GetEnv[string]("MYSQL_HOST"),
		utils.GetEnvWithDefault("MYSQL_PORT", 3306),
		utils.GetEnv[string]("MYSQL_DB"),
	)
}
