package pager_duty

import (
	"testing"
	"time"

	"github.com/h2non/gock"
)

// TestPagerDutySuccess tests the PagerDuty function for a successful API call.
func TestPagerDutySuccess(t *testing.T) {
	defer gock.Off()

	currentTime := time.Date(2024, 11, 17, 20, 34, 58, 651387237, time.UTC)
	startTime := currentTime.Add(-1 * time.Hour).Format(time.RFC3339)
	endTime := currentTime.Format(time.RFC3339)

	mockApiKey := "mock_apiKey"
	mockSchedule := "mock_schedule"
	mockUserEmail := "test@example.com"

	// Mock the GetEnvVars function
	originalGetEnvVars := GetEnvVars
	defer func() { GetEnvVars = originalGetEnvVars }()

	GetEnvVars = func() map[string]string {
		return map[string]string{
			"apiKey":   mockApiKey,
			"schedule": mockSchedule,
		}
	}

	gock.New("https://api.pagerduty.com").
		MatchHeader(headerAuthorization, "Token token="+mockApiKey).
		MatchHeader(headerAccept, acceptHeaderType).
		Get("/schedules/" + mockSchedule + "/users").
		MatchParams(map[string]string{"since": startTime, "until": endTime}).
		Reply(200).
		JSON(map[string][]User{
			"users": {{Name: "Test User", Email: mockUserEmail}},
		})

	result, err := PagerDuty(currentTime)
	if err != nil {
		t.Errorf("PagerDuty returned an error: %s", err)
	}
	expected := "On-Call SIRT " + mockUserEmail
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// TestPagerDutyNoUsers tests the scenario where the API returns an empty users array.
func TestPagerDutyNoUsers(t *testing.T) {
	defer gock.Off()

	currentTime := time.Date(2024, 11, 17, 20, 34, 58, 651387237, time.UTC)

	mockApiKey := "mock_apiKey"
	mockSchedule := "mock_schedule"

	originalGetEnvVars := GetEnvVars
	defer func() { GetEnvVars = originalGetEnvVars }()

	GetEnvVars = func() map[string]string {
		return map[string]string{
			"apiKey":   mockApiKey,
			"schedule": mockSchedule,
		}
	}

	startTime := currentTime.Add(-1 * time.Hour).Format(time.RFC3339)
	endTime := currentTime.Format(time.RFC3339)

	gock.New("https://api.pagerduty.com").
		Get("/schedules/"+mockSchedule+"/users").
		MatchParams(map[string]string{"since": startTime, "until": endTime}).
		MatchHeader(headerAuthorization, "Token token="+mockApiKey).
		MatchHeader(headerAccept, acceptHeaderType).
		Reply(200).
		JSON(map[string][]User{
			"users": {{}},
		})

	result, err := PagerDuty(currentTime)
	if err != nil {
		t.Errorf("PagerDuty returned an error: %s", err)
	}
	expected := "On-Call SIRT "
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

}

// TestPagerDutyApiError tests the PagerDuty function when the API request fails.
func TestPagerDutyApiError(t *testing.T) {
	defer gock.Off()

	currentTime := time.Date(2024, 11, 17, 20, 34, 58, 651387237, time.UTC)

	mockApiKey := "mock_apiKey"
	mockSchedule := "mock_schedule"

	originalGetEnvVars := GetEnvVars
	defer func() { GetEnvVars = originalGetEnvVars }()

	GetEnvVars = func() map[string]string {
		return map[string]string{
			"apiKey":   mockApiKey,
			"schedule": mockSchedule,
		}
	}
	startTime := currentTime.Add(-1 * time.Hour).Format(time.RFC3339)
	endTime := currentTime.Format(time.RFC3339)

	gock.New("https://api.pagerduty.com").
		Get("/schedules/"+mockSchedule+"/users").
		MatchParams(map[string]string{"since": startTime, "until": endTime}).
		MatchHeader(headerAuthorization, "Token token="+mockApiKey).
		Reply(400).
		JSON(map[string]string{"error": "Something went wrong"})

	_, err := PagerDuty(currentTime)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Check the panic message
	expectedError := "response status code is not 200. status code: 400"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

}

// TestPagerDutyJsonError tests the PagerDuty function when json doesn't parse.
func TestPagerDutyJsonError(t *testing.T) {
	defer gock.Off()

	currentTime := time.Date(2024, 11, 17, 20, 34, 58, 651387237, time.UTC)

	mockApiKey := "mock_apiKey"
	mockSchedule := "mock_schedule"

	originalGetEnvVars := GetEnvVars
	defer func() { GetEnvVars = originalGetEnvVars }()

	GetEnvVars = func() map[string]string {
		return map[string]string{
			"apiKey":   mockApiKey,
			"schedule": mockSchedule,
		}
	}
	startTime := currentTime.Add(-1 * time.Hour).Format(time.RFC3339)
	endTime := currentTime.Format(time.RFC3339)

	gock.New("https://api.pagerduty.com").
		Get("/schedules/"+mockSchedule+"/users").
		MatchParams(map[string]string{"since": startTime, "until": endTime}).
		MatchHeader(headerAuthorization, "Token token="+mockApiKey).
		Reply(200).
		BodyString("invalid json")

	_, err := PagerDuty(currentTime)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Check the panic message
	expectedError := "error parsing json: invalid character 'i' looking for beginning of value"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

}
