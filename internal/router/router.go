package router

import (
	"github.com/gin-gonic/gin"
	"github.com/meetgeekai/meeting_service/internal/handlers"
	"github.com/meetgeekai/meeting_service/internal/services/meetings"
	"go.uber.org/zap"
)

func RegisterRoutes(logger *zap.Logger, router *gin.RouterGroup, meetingsSvc *meetings.MeetingsService) {
	meetingsHandler := handlers.NewMeetingsHandler(meetingsSvc, logger)

	router.GET("/upcoming-meetings", meetingsHandler.GetUpcomingMeetings)
	router.PATCH("/update-auto-join", meetingsHandler.UpdateAutoJoin)
}
