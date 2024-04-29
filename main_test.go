package main

import (
	"testing"
	"time"
)

func Test_getLastWeek(t *testing.T) {
	tests := []struct {
		name     string
		today    time.Time
		lastweek time.Time
	}{
		{"success", time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC)},
	}

	nowFunc = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2024-03-01T08:30:00Z")
		return t
	}
	defer func() {
		nowFunc = nil
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			today, lw := getLastWeek()
			if (today != tt.today) || (lw != tt.lastweek) {
				t.Errorf("expected) today: %s, lastweek: %s", tt.today, tt.lastweek)
				t.Errorf("  actual) today: %s, lastweek: %s", today, lw)
			}
		})
	}
}
