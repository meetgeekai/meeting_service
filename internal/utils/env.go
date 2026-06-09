package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/meetgeekai/go-common/logging"
	"github.com/meetgeekai/go-common/utils"
	"go.uber.org/zap"
)

type Environment string

const (
	Local       Environment = "local"
	Development Environment = "development"
	Production  Environment = "production"
)

func GetEnvironment() Environment {
	env := utils.GetEnvWithDefault("ENVIRONMENT", "local")

	switch env {
	case "local":
		return Local
	case "dev", "development":
		return Development
	case "prod", "production":
		return Production
	default:
		panic("unknown environment")
	}
}

func IsLocal() bool {
	return GetEnvironment() == Local
}

func IsDevelopment() bool {
	return GetEnvironment() == Development
}

func IsProduction() bool {
	return GetEnvironment() == Production
}

func GetGinMode() string {
	switch GetEnvironment() {
	case Local, Development:
		return gin.DebugMode
	case Production:
		return gin.ReleaseMode
	default:
		panic("unknown environment")
	}
}

func GetLogger() *zap.Logger {
	logsDir := utils.GetEnv[string]("LOGS_DIR")
	slackLogLevel := utils.GetEnv[string]("SLACK_LOG_LEVEL")
	slackWebhookUrl := utils.GetEnv[string]("SLACK_WEBHOOK_URL")

	switch GetEnvironment() {
	case Local:
		return logging.CreateLocalLogger()
	case Development:
		return logging.CreateDevelopmentLogger(logsDir)
	case Production:
		return logging.CreateProductionLogger(logsDir, slackLogLevel, slackWebhookUrl)
	default:
		panic("unknown environment")
	}
}
