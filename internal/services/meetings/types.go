package meetings

import (
	"github.com/meetgeekai/meeting_service/internal/models"
	"github.com/meetgeekai/meeting_service/internal/repositories"
	elasticsearch "github.com/meetgeekai/meeting_service/internal/services/es"
	"go.uber.org/zap"
)

const (
	pageSize            = 30
	pitKeepAliveSeconds = 30
)

type MeetingsService struct {
	repo   repositories.MeetingsRepository
	es     *elasticsearch.ESService
	logger *zap.Logger
}

func NewMeetingsService(
	repo repositories.MeetingsRepository,
	es *elasticsearch.ESService,
	logger *zap.Logger,
) *MeetingsService {
	return &MeetingsService{
		repo:   repo,
		es:     es,
		logger: logger,
	}
}

type BotStatus struct {
	Status    *string `json:"status"`
	SubStatus *string `json:"sub_status"`
	SessionID *string `json:"session_id"`
}

type ConversationTemplate struct {
	Name *string `json:"name"`
	ID   *int64  `json:"id"`
}

type UpcomingMeeting struct {
	OwnerUUID            string                        `json:"owner_uuid"`
	OwnerName            string                        `json:"owner_name"`
	OwnerEmail           string                        `json:"owner_email"`
	HostEmail            *string                       `json:"host_email"`
	StartTime            *string                       `json:"start_time"`
	EndTime              *string                       `json:"end_time"`
	StartTimeActual      *string                       `json:"start_time_actual"`
	EndTimeActual        *string                       `json:"end_time_actual"`
	ConversationTemplate ConversationTemplate          `json:"conversation_template"`
	Source               *string                       `json:"source"`
	Language             *models.TranscriptionLanguage `json:"language"`
	MeetingID            *string                       `json:"meeting_id"`
	ConferenceVendor     *string                       `json:"conference_vendor"`
	Title                *string                       `json:"title"`
	EditableTitle        *string                       `json:"editable_title"`
	IDES                 *string                       `json:"id_es"`
	ParticipantsActual   []any                         `json:"participants_actual"`
	Participants         []any                         `json:"participants"`
	Agenda               string                        `json:"agenda"`
	EditableAgenda       *string                       `json:"editable_agenda"`
	ExternalRecurringIDs []any                         `json:"external_recurring_ids"`
	VoiceAgentID         *int64                        `json:"voice_agent_id"`
	DynamicTemplate      bool                          `json:"dynamic_template"`
	AccessType           *string                       `json:"access_type"`
	AutomaticJoin        bool                          `json:"automatic_join"`
	BotStatus            *BotStatus                    `json:"bot_status"`
	Analyzing            *bool                         `json:"analyzing"`
	JoinURL              string                        `json:"join_url"`
}

type UpcomingMeetingsPage struct {
	Meetings   []UpcomingMeeting `json:"meetings"`
	NextCursor *string           `json:"next_cursor"`
}

type UpdateAutoJoinResult struct {
	MeetingID     string `json:"meeting_id"`
	AutomaticJoin bool   `json:"automatic_join"`
}
