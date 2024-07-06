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
		{"success", time.Date(2024, 3, 1, 0, 0, 0, 0, time.Local), time.Date(2024, 2, 23, 0, 0, 0, 0, time.Local)},
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

func Test_getAggregationPeriod(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	today, lastweek := getLastWeek()

	type args struct {
		sinceStr string
		untilStr string
	}
	tests := []struct {
		name    string
		args    args
		since   time.Time
		until   time.Time
		wantErr bool
	}{
		{"success: arguments are UTC", args{sinceStr: "2024-06-01T00:00:00Z", untilStr: "2024-07-01T00:00:00Z"}, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), false},
		{"success: arguments are other TZ", args{sinceStr: "2024-06-01T00:00:00+09:00", untilStr: "2024-07-01T00:00:00+09:00"}, time.Date(2024, 6, 1, 0, 0, 0, 0, jst), time.Date(2024, 7, 1, 0, 0, 0, 0, jst), false},
		{"success: arguments are empty string", args{sinceStr: "", untilStr: ""}, lastweek, today, false},
		{"failure: fist argument format is invalid", args{sinceStr: "2024-06-01T00:00:00", untilStr: "2024-07-01T00:00:00Z"}, time.Time{}, time.Time{}, true},
		{"failure: second argument format is invalid", args{sinceStr: "2024-06-01T00:00:00Z", untilStr: "2024-07-01T00:00:00"}, time.Time{}, time.Time{}, true},
		{"failure: since is after until", args{sinceStr: "2024-06-01T00:00:00+00:00", untilStr: "2024-05-01T00:00:00+00:00"}, time.Time{}, time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, err := getAggregationPeriod(tt.args.sinceStr, tt.args.untilStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAggregationPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !got1.Equal(tt.since) {
				t.Errorf("getAggregationPeriod() got1 = %v, want %v", got1, tt.since)
			}
			if !got2.Equal(tt.until) {
				t.Errorf("getAggregationPeriod() got2 = %v, want %v", got2, tt.until)
			}
		})
	}
}
