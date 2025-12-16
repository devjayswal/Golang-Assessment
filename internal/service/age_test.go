package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name string
		dob  time.Time
		now  time.Time
		want int
	}{
		{"before_birthday", time.Date(2000, 12, 10, 0, 0, 0, 0, time.UTC), time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC), 23},
		{"on_birthday", time.Date(2000, 12, 10, 0, 0, 0, 0, time.UTC), time.Date(2024, 12, 10, 0, 0, 0, 0, time.UTC), 24},
		{"after_birthday", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), 24},
		{"future_date", time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateAge(tt.dob, tt.now); got != tt.want {
				t.Fatalf("want %d got %d", tt.want, got)
			}
		})
	}
}
