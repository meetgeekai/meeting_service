package meetings

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	elasticsearch "github.com/meetgeekai/meeting_service/internal/services/es"
)

func encodeCursor(c elasticsearch.UpcomingMeetingsCursor) (string, error) {
	payload, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func decodeCursor(encoded string) (*elasticsearch.UpcomingMeetingsCursor, error) {
	payload, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}
	var c elasticsearch.UpcomingMeetingsCursor
	if err := json.Unmarshal(payload, &c); err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}
	if c.PitID == "" || len(c.SearchAfter) == 0 || bytes.Equal(c.SearchAfter, []byte("null")) {
		return nil, fmt.Errorf("invalid cursor: missing pit_id or search_after")
	}
	return &c, nil
}
