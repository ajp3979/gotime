package pager_duty

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

const apiKeyEnvVar = "PAGER_DUTY_API_KEY"
const scheduleEnvVar = "PAGER_DUTY_SCHEDULE"
const headerAuthorization = "Authorization"
const headerAccept = "Accept"
const acceptHeaderType = "application/vnd.pagerduty+json;version=2"

var GetEnvVars = getEnvVars

func getEnvVars() map[string]string {
	apiKey := os.Getenv(apiKeyEnvVar)
	if apiKey == "" {
		log.Fatal("API key env not set")
	}
	schedule := os.Getenv(scheduleEnvVar)
	if schedule == "" {
		log.Fatal("Schedule env not set")
	}
	return map[string]string{
		"apiKey":   apiKey,
		"schedule": schedule,
	}
}

func PagerDuty(currentTime time.Time) (string, error) {
	enVars := GetEnvVars()
	startTime := currentTime.Add(-1 * time.Hour).Format(time.RFC3339)
	endTime := currentTime.Format(time.RFC3339)
	url := "https://api.pagerduty.com/schedules/" + enVars["schedule"] + "/users?since=" + startTime + "&until=" + endTime

	headers := map[string]string{
		headerAuthorization: "Token token=" + enVars["apiKey"],
		headerAccept:        acceptHeaderType,
	}

	res, err := MakeHTTPRequest(url, headers)
	if err != nil {
		error := fmt.Errorf("error making http request: %s", err)
		return "", error
	}

	if res.StatusCode != http.StatusOK {
		error := fmt.Errorf("response status code is not 200. status code: %v", res.StatusCode)
		return "", error
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Failed to close body: %s", err)
		}
	}(res.Body)

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		error := fmt.Errorf("failed to read response body: %s", readErr)
		return "", error
	}

	var responseData map[string][]User
	if err := json.Unmarshal(body, &responseData); err != nil {
		error := fmt.Errorf("error parsing json: %s", err)
		return "", error
	}

	users := responseData["users"]
	return "On-Call SIRT " + users[0].Email, nil
}

func MakeHTTPRequest(url string, headers map[string]string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	return client.Do(req)
}
