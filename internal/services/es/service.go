package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	backoff "github.com/cenkalti/backoff/v4"
	commonBackoff "github.com/meetgeekai/go-common/backoff"
	"go.uber.org/zap"
)

func (s *ESService) callWithResponse(ctx context.Context, path string, body any, logFields ...zap.Field) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("%s%s", s.config.ESBase, path)

	var respBody []byte
	var status int
	err = commonBackoff.WithRetryBackoff(commonBackoff.RetryBackoffConfig{
		Algorithm: backoff.NewExponentialBackOff(),
		MaxTries:  5,
		Func: func() error {
			req, err := http.NewRequestWithContext(ctx,
				"GET", // Ugh.... It is what it is.
				url,
				bytes.NewBuffer(payload))
			if err != nil {
				return backoff.Permanent(err)
			}

			req.Header.Add("Authorization", s.config.APISecret)

			resp, err := s.client.Do(req)
			if err != nil {
				if ctx.Err() != nil {
					return backoff.Permanent(ctx.Err())
				}
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 500 {
				return fmt.Errorf("ES service returned %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			respBody = data
			status = resp.StatusCode
			return nil
		},
		OnRetry: func(err error) {
			s.logger.Warn("Retrying ES call", append([]zap.Field{zap.Error(err), zap.String("path", path)}, logFields...)...)
		},
		OnGiveup: func(err error) {
			s.logger.Error("Gave up ES call", append([]zap.Field{zap.Error(err), zap.String("path", path)}, logFields...)...)
		},
	})

	if err != nil {
		return nil, err
	}

	if status >= 400 {
		return nil, fmt.Errorf("ES service returned %d", status)
	}

	return respBody, nil
}

type UpcomingMeetingsCursor struct {
	PitID       string          `json:"pit_id"`
	SearchAfter json.RawMessage `json:"search_after"`
}

type UpcomingMeetingsPageResult struct {
	NextCursor *UpcomingMeetingsCursor
	Meetings   []json.RawMessage
}

func (s *ESService) GetUpcomingMeetingsPage(
	ctx context.Context,
	ownerUUID string,
	cursor *UpcomingMeetingsCursor,
	size int,
	pitKeepAliveSeconds int,
) (*UpcomingMeetingsPageResult, error) {
	body := map[string]any{
		"owner_uuid":             ownerUUID,
		"size":                   size,
		"pit_keep_alive_seconds": pitKeepAliveSeconds,
	}
	if cursor != nil {
		body["pit_id"] = cursor.PitID
		body["search_after"] = cursor.SearchAfter
	}

	raw, err := s.callWithResponse(
		ctx,
		s.config.ESGetUpcomingMeetingsPage,
		body,
		zap.String("ownerUUID", ownerUUID),
	)
	if err != nil {
		return nil, err
	}

	var tuple [3]json.RawMessage
	if err := json.Unmarshal(raw, &tuple); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ES response: %w", err)
	}

	var nextPitID *string
	if err := json.Unmarshal(tuple[0], &nextPitID); err != nil {
		return nil, fmt.Errorf("failed to parse next_pit_id: %w", err)
	}

	var nextSearchAfter []any
	if err := json.Unmarshal(tuple[1], &nextSearchAfter); err != nil {
		return nil, fmt.Errorf("failed to parse next_search_after: %w", err)
	}

	var meetings []json.RawMessage
	if err := json.Unmarshal(tuple[2], &meetings); err != nil {
		return nil, fmt.Errorf("failed to parse meetings: %w", err)
	}

	result := &UpcomingMeetingsPageResult{Meetings: meetings}
	if nextPitID != nil && len(nextSearchAfter) > 0 {
		result.NextCursor = &UpcomingMeetingsCursor{
			PitID:       *nextPitID,
			SearchAfter: tuple[1],
		}
	}

	return result, nil
}
