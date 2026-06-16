package repositories

import (
	"context"

	models "github.com/meetgeekai/meeting_service/internal/models"
)

type MeetingsRepository interface {
	GetUserForUpcomingMeetings(ctx context.Context, userUUID string) (*models.UpcomingMeetingsOwner, error)
	GetAvailableTranscriptionLanguages(ctx context.Context) ([]models.TranscriptionLanguage, error)
	GetConversationTemplateNames(ctx context.Context, ids []int64, userID uint32) (map[int64]string, error)
	GetConnectedCalendarVendors(ctx context.Context, userUUID string) (models.ConnectedCalendars, error)
}
