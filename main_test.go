package main

import "testing"

func TestClock(t *testing.T) {
	if got := clock("2026-06-26T14:30"); got != "14:30" {
		t.Fatalf("clock() = %q, want %q", got, "14:30")
	}
}

func TestHourlyStart(t *testing.T) {
	hourly := Hourly{
		Time: []string{
			"2026-06-26T12:00",
			"2026-06-26T13:00",
			"2026-06-26T14:00",
		},
	}

	if got := hourlyStart(hourly, "2026-06-26T13:15"); got != 2 {
		t.Fatalf("hourlyStart() = %d, want %d", got, 2)
	}
}

func TestWeather(t *testing.T) {
	if got := weather(0); got != "Clear sky" {
		t.Fatalf("weather() = %q, want %q", got, "Clear sky")
	}
}
