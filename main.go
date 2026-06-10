package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	hs "github.com/meetgeekai/go-common/healthserver"
	gocommon "github.com/meetgeekai/go-common/utils"
	"github.com/meetgeekai/meeting_service/internal/middleware"
	mysqlRepo "github.com/meetgeekai/meeting_service/internal/repositories/mysql"
	"github.com/meetgeekai/meeting_service/internal/router"
	elasticsearch "github.com/meetgeekai/meeting_service/internal/services/es"
	"github.com/meetgeekai/meeting_service/internal/services/meetings"
	"github.com/meetgeekai/meeting_service/internal/utils"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
)

var logger *zap.Logger
var hsv *hs.HealthServer

func init() {
	logger = utils.GetLogger()
	hsv = hs.NewHealthServer(logger).Start()
	logger.Error("Service :: Server restarted")
}

func main() {
	defer logger.Sync()

	gin.SetMode(utils.GetGinMode())
	engine := gin.Default()

	srv := &http.Server{
		Addr:    utils.GetAddr(),
		Handler: engine,
	}

	limiter := ratelimit.New(gocommon.GetEnv[int]("GLOBAL_RATE_LIMIT"))
	secret := gocommon.GetEnv[string]("GLOBAL_API_SECRET")

	repo := mysqlRepo.NewMySQLRepository()
	esService := elasticsearch.NewESService(&elasticsearch.ESServiceConfig{
		ESBase:                             gocommon.GetEnv[string]("ES_SERVICE_ENDPOINT"),
		ESGetUpcomingMeetingsPage:          gocommon.GetEnv[string]("ES_SERVICE_GET_UPCOMING_MEETINGS_PAGE"),
		ESGetUpcomingMeetingByID:           gocommon.GetEnv[string]("ES_SERVICE_GET_UPCOMING_MEETING_BY_ID"),
		ESUpdateMeetingPartially:           gocommon.GetEnv[string]("ES_SERVICE_UPDATE_MEETING_PARTIALLY"),
		ESUpdateRecurrentMeetingsPartially: gocommon.GetEnv[string]("ES_SERVICE_UPDATE_RECURRENT_MEETINGS_PARTIALLY"),
		APISecret:                          secret,
	}, logger)
	meetingsService := meetings.NewMeetingsService(repo, esService, logger)

	routerGroup := engine.Group("/")
	routerGroup.Use(middleware.RateLimiterMiddleware(limiter))
	routerGroup.Use(middleware.EntryMiddleware())
	routerGroup.Use(middleware.AuthMiddleware(secret))

	router.RegisterRoutes(logger, routerGroup, meetingsService)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("GIN Server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Microservice received SIGTERM signal. It's now scheduled to shutdown in 30 seconds...")

	time.Sleep(2 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	if err := hsv.Shutdown(ctx); err != nil {
		logger.Fatal("Health Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
