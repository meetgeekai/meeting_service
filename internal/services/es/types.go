package elasticsearch

import (
	"net/http"

	"go.uber.org/zap"
)

type ESServiceConfig struct {
	ESBase                             string
	ESGetUpcomingMeetingsPage          string
	ESGetUpcomingMeetingByID           string
	ESUpdateMeetingPartially           string
	ESUpdateRecurrentMeetingsPartially string
	APISecret                          string
}

type ESService struct {
	config *ESServiceConfig
	client http.Client
	logger *zap.Logger
}

func NewESService(config *ESServiceConfig, logger *zap.Logger) *ESService {
	return &ESService{
		config: config,
		client: http.Client{},
		logger: logger,
	}
}
