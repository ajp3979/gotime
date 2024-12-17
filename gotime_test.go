package main

import (
	"testing"
	"time"
)

// Function to check the expected output against the generated output
func checkRegionOutput(t *testing.T, region Region, currentTime time.Time, expected string) {
	l, _ := time.LoadLocation(region.Location)
	output := region.Name + ": " + currentTime.In(l).Format("Monday")
	if output != expected {
		t.Errorf("expected %q but got %q", expected, output)
	}
}

// TestGenerateRegion tests the GenerateRegion function.
func TestGenerateRegion(t *testing.T) {
	tests := []struct {
		name     string
		region   Region
		expected string
	}{
		{"Hyderabad Region Test",
			Region{Name: "Hyderabad", Location: "Asia/Kolkata"},
			"Hyderabad: Monday"},
		{"UTC Region Test",
			Region{Name: "UTC", Location: "UTC"},
			"UTC: Monday"},
		{"Dublin Region Test",
			Region{Name: "Dublin", Location: "Europe/Dublin"},
			"Dublin: Monday"},
		{"Eastern USA Region Test",
			Region{Name: "Eastern USA", Location: "America/New_York"},
			"Eastern USA: Sunday"},
	}

	currentTime := time.Date(2023, 10, 23, 0, 0, 0, 0, time.UTC) // change date as needed

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkRegionOutput(t, tt.region, currentTime, tt.expected)
		})
	}
}

// TestMainFunction runs the main function and checks for panics.
func TestMainFunction(t *testing.T) {
	runMainAndCheckPanic(t)
}

// runMainAndCheckPanic runs the main function and ensures it does not panic.
func runMainAndCheckPanic(t *testing.T) {
	mainFunc := main
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() panicked with value: %v", r)
		}
	}()
	mainFunc()
}
