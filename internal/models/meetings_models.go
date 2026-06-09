package models

import (
	"strings"
)

type UpcomingMeetingsOwner struct {
	UUID  string
	Name  string
	Email string
}

// CalendarProvider identifies a calendar provider. It is matched case-insensitively against both
// the meeting `source` stored in Elasticsearch and the integration `vendor` stored in MySQL.
type CalendarProvider string

const (
	CalendarProviderGoogle    CalendarProvider = "google"
	CalendarProviderMicrosoft CalendarProvider = "microsoft"
)

// ParseCalendarProvider maps an arbitrary source/vendor string to a known CalendarProvider,
// reporting false when it is neither Google nor Microsoft.
func ParseCalendarProvider(value string) (CalendarProvider, bool) {
	switch {
	case strings.EqualFold(value, string(CalendarProviderGoogle)):
		return CalendarProviderGoogle, true
	case strings.EqualFold(value, string(CalendarProviderMicrosoft)):
		return CalendarProviderMicrosoft, true
	default:
		return "", false
	}
}

// ConnectedCalendars indicates which calendar providers a user has actively connected.
type ConnectedCalendars struct {
	Google    bool
	Microsoft bool
}

// Allows reports whether a meeting with the given source should be kept, mirroring the PHP
// filterMeetingsByConnectedCalendars: calendar-sourced meetings are kept only when the matching
// calendar is connected; meetings from any other source are always kept.
func (c ConnectedCalendars) Allows(source *string) bool {
	if source == nil {
		return true
	}
	provider, ok := ParseCalendarProvider(*source)
	if !ok {
		return true
	}
	switch provider {
	case CalendarProviderGoogle:
		return c.Google
	case CalendarProviderMicrosoft:
		return c.Microsoft
	default:
		return true
	}
}

type TranscriptionLanguage struct {
	ID               int64  `json:"id"`
	Code             string `json:"code"`
	Value            string `json:"value"`
	Country          string `json:"country"`
	Language         string `json:"language"`
	CustomDictionary bool   `json:"custom_dictionary"`
}
