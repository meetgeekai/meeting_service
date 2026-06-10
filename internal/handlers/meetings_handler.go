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
		UserID uint32 `form:"user_id" binding:"required"`
		Cursor string `form:"cursor"`
	}

	var req Request
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AppError{
			Status: response.INVALID_INPUT_ERR,
			Reason: err.Error(),
		})
		return
	}

	resp := h.svc.GetUpcomingMeetings(c, req.UserID, req.Cursor)
	if resp.IsError() {
		c.JSON(resp.GetHttpCode(), resp.Err)
		return
	}

	c.JSON(resp.GetHttpCode(), resp.Data)
}

func (h *MeetingsHandler) UpdateAutoJoin(c *gin.Context) {
	type Request struct {
		MeetingID     string `json:"meeting_id" binding:"required"`
		AutomaticJoin *bool  `json:"automatic_join" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AppError{
			Status: response.INVALID_INPUT_ERR,
			Reason: err.Error(),
		})
		return
	}

	resp := h.svc.UpdateAutoJoin(c, req.MeetingID, *req.AutomaticJoin)
	if resp.IsError() {
		c.JSON(resp.GetHttpCode(), resp.Err)
		return
	}

	c.JSON(resp.GetHttpCode(), resp.Data)
}
