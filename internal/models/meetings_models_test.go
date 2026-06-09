package models

import "testing"

func TestConnectedCalendarsAllows(t *testing.T) {
	strPtr := func(s string) *string { return &s }

	cases := []struct {
		name      string
		calendars ConnectedCalendars
		source    *string
		want      bool
	}{
		{"nil source always kept", ConnectedCalendars{}, nil, true},
		{"unknown source always kept", ConnectedCalendars{}, strPtr("adhoc"), true},
		{"google kept when connected", ConnectedCalendars{Google: true}, strPtr("google"), true},
		{"google dropped when not connected", ConnectedCalendars{Microsoft: true}, strPtr("google"), false},
		{"microsoft kept when connected", ConnectedCalendars{Microsoft: true}, strPtr("microsoft"), true},
		{"microsoft dropped when not connected", ConnectedCalendars{Google: true}, strPtr("microsoft"), false},
		{"google matched case-insensitively", ConnectedCalendars{Google: true}, strPtr("Google"), true},
		{"microsoft matched case-insensitively", ConnectedCalendars{Microsoft: true}, strPtr("MICROSOFT"), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.calendars.Allows(tc.source); got != tc.want {
				t.Errorf("Allows(%v) = %v, want %v", tc.source, got, tc.want)
			}
		})
	}
}
