package meetings

import (
	"context"
	"encoding/json"

	"github.com/meetgeekai/meeting_service/internal/models"
	"github.com/meetgeekai/meeting_service/internal/response"
	elasticsearch "github.com/meetgeekai/meeting_service/internal/services/es"
	"go.uber.org/zap"
)

type esMeeting struct {
	HostEmail               *string    `json:"host_email"`
	StartTime               *string    `json:"start_time"`
	EndTime                 *string    `json:"end_time"`
	StartTimeActual         *string    `json:"start_time_actual"`
	EndTimeActual           *string    `json:"end_time_actual"`
	Source                  *string    `json:"source"`
	Language                *string    `json:"language"`
	MeetingID               *string    `json:"meeting_id"`
	ConferenceVendor        *string    `json:"conference_vendor"`
	Title                   *string    `json:"title"`
	EditableTitle           *string    `json:"editable_title"`
	IDES                    *string    `json:"es_id"`
	ParticipantsActual      []any      `json:"participants_actual"`
	Participants            []any      `json:"participants"`
	Agenda                  string     `json:"agenda"`
	EditableAgenda          *string    `json:"editable_agenda"`
	ExternalRecurringIDs    []any      `json:"external_recurring_ids"`
	VoiceAgentID            *int64     `json:"voice_agent_id"`
	DynamicTemplate         bool       `json:"dynamic_template"`
	AccessType              *string    `json:"access_type"`
	AutomaticJoin           bool       `json:"automatic_join"`
	BotStatus               *BotStatus `json:"bot_status"`
	Analyzing               *bool      `json:"analyzing"`
	JoinURL                 string     `json:"join_link"`
	ConversationTemplateIDs []int64    `json:"conversation_template_ids"`
}

func (s *MeetingsService) GetUpcomingMeetings(
	ctx context.Context,
	userUUID string,
	rawCursor string,
) response.Response[UpcomingMeetingsPage] {
	owner, err := s.repo.GetUserForUpcomingMeetings(ctx, userUUID)
	if err != nil {
		s.logger.Error("Failed to fetch user", zap.String("userUUID", userUUID), zap.Error(err))
		return response.Error[UpcomingMeetingsPage](response.DB_ERR, "Failed to fetch user")
	}
	if owner == nil {
		return response.Error[UpcomingMeetingsPage](response.NOT_FOUND_ERR, "User %s not found", userUUID)
	}

	var cursor *elasticsearch.UpcomingMeetingsCursor
	if rawCursor != "" {
		cursor, err = decodeCursor(rawCursor)
		if err != nil {
			return response.Error[UpcomingMeetingsPage](response.INVALID_INPUT_ERR, "%s", err.Error())
		}
	}

	page, err := s.es.GetUpcomingMeetingsPage(ctx, owner.UUID, cursor, pageSize, pitKeepAliveSeconds)
	if err != nil {
		s.logger.Error("Failed to fetch upcoming meetings page",
			zap.String("ownerUUID", owner.UUID), zap.Error(err))
		return response.Error[UpcomingMeetingsPage](response.SERVICE_ERR, "Failed to fetch upcoming meetings")
	}

	calendars, err := s.repo.GetConnectedCalendarVendors(ctx, owner.UUID)
	if err != nil {
		s.logger.Error("Failed to fetch connected calendars",
			zap.String("ownerUUID", owner.UUID), zap.Error(err))
		return response.Error[UpcomingMeetingsPage](response.DB_ERR, "Failed to fetch connected calendars")
	}

	languages, err := s.repo.GetAvailableTranscriptionLanguages(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch transcription languages", zap.Error(err))
		return response.Error[UpcomingMeetingsPage](response.DB_ERR, "Failed to fetch transcription languages")
	}

	languagesByCode := make(map[string]models.TranscriptionLanguage, len(languages))
	for _, lang := range languages {
		languagesByCode[lang.Code] = lang
	}

	rawMeetings := make([]esMeeting, 0, len(page.Meetings))
	templateIDSet := make(map[int64]struct{})
	for _, rawMsg := range page.Meetings {
		var raw esMeeting
		if err := json.Unmarshal(rawMsg, &raw); err != nil {
			s.logger.Error("Failed to unmarshal meeting from ES", zap.Error(err))
			continue
		}
		if !calendars.Allows(raw.Source) {
			continue
		}
		rawMeetings = append(rawMeetings, raw)
		if len(raw.ConversationTemplateIDs) > 0 {
			templateIDSet[raw.ConversationTemplateIDs[0]] = struct{}{}
		}
	}

	var templateNames map[int64]string
	if len(templateIDSet) > 0 {
		templateIDs := make([]int64, 0, len(templateIDSet))
		for id := range templateIDSet {
			templateIDs = append(templateIDs, id)
		}
		templateNames, err = s.repo.GetConversationTemplateNames(ctx, templateIDs, owner.ID)
		if err != nil {
			s.logger.Error("Failed to fetch conversation template names", zap.Error(err))
			return response.Error[UpcomingMeetingsPage](response.DB_ERR, "Failed to fetch conversation template names")
		}
	}

	meetings := make([]UpcomingMeeting, 0, len(rawMeetings))
	for _, raw := range rawMeetings {
		meetings = append(meetings, transformMeeting(raw, owner, languagesByCode, templateNames))
	}

	result := UpcomingMeetingsPage{Meetings: meetings}

	if page.NextCursor != nil {
		encoded, err := encodeCursor(*page.NextCursor)
		if err != nil {
			s.logger.Error("Failed to encode cursor", zap.Error(err))
			return response.Error[UpcomingMeetingsPage](response.SERVICE_ERR, "Failed to encode cursor")
		}
		result.NextCursor = &encoded
	}

	return response.Success(result)
}

func transformMeeting(
	raw esMeeting,
	owner *models.UpcomingMeetingsOwner,
	languagesByCode map[string]models.TranscriptionLanguage,
	templateNames map[int64]string,
) UpcomingMeeting {
	participantsActual := raw.ParticipantsActual
	if participantsActual == nil {
		participantsActual = []any{}
	}
	participants := raw.Participants
	if participants == nil {
		participants = []any{}
	}

	meeting := UpcomingMeeting{
		OwnerUUID:            owner.UUID,
		OwnerName:            owner.Name,
		OwnerEmail:           owner.Email,
		HostEmail:            raw.HostEmail,
		StartTime:            raw.StartTime,
		EndTime:              raw.EndTime,
		StartTimeActual:      raw.StartTimeActual,
		EndTimeActual:        raw.EndTimeActual,
		Source:               raw.Source,
		MeetingID:            raw.MeetingID,
		ConferenceVendor:     raw.ConferenceVendor,
		Title:                raw.Title,
		EditableTitle:        raw.EditableTitle,
		IDES:                 raw.IDES,
		ParticipantsActual:   participantsActual,
		Participants:         participants,
		Agenda:               raw.Agenda,
		EditableAgenda:       raw.EditableAgenda,
		ExternalRecurringIDs: raw.ExternalRecurringIDs,
		VoiceAgentID:         raw.VoiceAgentID,
		DynamicTemplate:      raw.DynamicTemplate,
		AccessType:           raw.AccessType,
		AutomaticJoin:        raw.AutomaticJoin,
		BotStatus:            raw.BotStatus,
		Analyzing:            raw.Analyzing,
		JoinURL:              raw.JoinURL,
		ConversationTemplate: ConversationTemplate{},
	}

	if len(raw.ConversationTemplateIDs) > 0 {
		id := raw.ConversationTemplateIDs[0]
		meeting.ConversationTemplate.ID = &id
		if name, ok := templateNames[id]; ok {
			meeting.ConversationTemplate.Name = &name
		}
	}

	if raw.Language != nil {
		if lang, found := languagesByCode[*raw.Language]; found {
			meeting.Language = &lang
		}
	}

	return meeting
}
