package meetings

import (
	"context"

	"github.com/meetgeekai/meeting_service/internal/response"
	"go.uber.org/zap"
)

func (s *MeetingsService) UpdateAutoJoin(
	ctx context.Context,
	meetingID string,
	automaticJoin bool,
) response.Response[UpdateAutoJoinResult] {
	meeting, err := s.es.GetUpcomingMeetingByID(ctx, meetingID)
	if err != nil {
		s.logger.Error("Failed to fetch upcoming meeting",
			zap.String("meetingID", meetingID), zap.Error(err))
		return response.Error[UpdateAutoJoinResult](response.SERVICE_ERR, "Failed to fetch upcoming meeting")
	}
	if meeting == nil {
		return response.Error[UpdateAutoJoinResult](response.NOT_FOUND_ERR, "Upcoming meeting %s not found", meetingID)
	}

	fields := map[string]any{"automatic_join": automaticJoin}

	updated, err := s.es.UpdateMeetingPartially(ctx, meeting.IDES, fields)
	if err != nil || !updated {
		s.logger.Error("Failed to update meeting",
			zap.String("meetingID", meetingID), zap.String("esID", meeting.IDES), zap.Error(err))
		return response.Error[UpdateAutoJoinResult](response.SERVICE_ERR, "Failed to update meeting %s", meetingID)
	}

	if len(meeting.ExternalRecurringIDs) > 0 {
		updated, err = s.es.UpdateRecurrentMeetingsPartially(
			ctx, meeting.OwnerUUID, meeting.ExternalRecurringIDs, meeting.IDES, fields)
		if err != nil || !updated {
			s.logger.Error("Failed to update recurrent meetings",
				zap.String("meetingID", meetingID), zap.String("ownerUUID", meeting.OwnerUUID), zap.Error(err))
			return response.Error[UpdateAutoJoinResult](response.SERVICE_ERR, "Failed to update recurrent meetings for meeting %s", meetingID)
		}
	}

	return response.Success(UpdateAutoJoinResult{
		MeetingID:     meetingID,
		AutomaticJoin: automaticJoin,
	})
}
