package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/meetgeekai/meeting_service/internal/response"
	"github.com/meetgeekai/meeting_service/internal/services/meetings"
)

type MeetingsHandler struct {
	svc    *meetings.MeetingsService
	logger *zap.Logger
}

func NewMeetingsHandler(svc *meetings.MeetingsService, logger *zap.Logger) *MeetingsHandler {
	return &MeetingsHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h *MeetingsHandler) GetUpcomingMeetings(c *gin.Context) {
	type Request struct {
		UserUUID string `form:"user_uuid" binding:"required"`
		Cursor   string `form:"cursor"`
	}

	var req Request
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AppError{
			Status: response.INVALID_INPUT_ERR,
			Reason: err.Error(),
		})
		return
	}

	resp := h.svc.GetUpcomingMeetings(c, req.UserUUID, req.Cursor)
	if resp.IsError() {
		c.JSON(resp.GetHttpCode(), resp.Err)
		return
	}

	c.JSON(resp.GetHttpCode(), resp.Data)
}

func (h *MeetingsHandler) UpdateAutoJoin(c *gin.Context) {
	meetingID := c.Param("meeting_id")
	if meetingID == "" {
		c.JSON(http.StatusBadRequest, response.AppError{
			Status: response.INVALID_INPUT_ERR,
			Reason: "meeting_id is required",
		})
		return
	}

	type Query struct {
		UserUUID string `form:"user_uuid" binding:"required"`
	}
	var query Query
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.AppError{
			Status: response.INVALID_INPUT_ERR,
			Reason: err.Error(),
		})
		return
	}

	type Body struct {
		AutomaticJoin *bool `json:"automatic_join" binding:"required"`
	}
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.AppError{
			Status: response.INVALID_INPUT_ERR,
			Reason: err.Error(),
		})
		return
	}

	resp := h.svc.UpdateAutoJoin(c, query.UserUUID, meetingID, *body.AutomaticJoin)
	if resp.IsError() {
		c.JSON(resp.GetHttpCode(), resp.Err)
		return
	}

	c.JSON(resp.GetHttpCode(), resp.Data)
}
