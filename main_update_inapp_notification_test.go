package NotificationAPI

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateInAppNotificationPassesWith202(t *testing.T) {
	Init(client_id, client_secret)
	params := InAppNotificationPatchRequest{
		TrackingIds: []string{"track123"},
		Opened:      stringPtr("2024-07-16T08:00:00Z"),
		Clicked:     stringPtr("2024-07-16T08:05:00Z"),
		Archived:    stringPtr("2024-07-16T09:00:00Z"),
		Actioned1:   stringPtr("Action1"),
		Actioned2:   stringPtr("Action2"),
		Reply: &struct {
			Date    string `json:"date"`
			Message string `json:"message"`
		}{
			Date:    "2024-07-16T10:00:00Z",
			Message: "Test reply message",
		},
	}
	userId := "123"
	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PATCH", "https://api.notificationapi.com/client_id/users/"+userId+"/notifications/INAPP_WEB",
		func(req *http.Request) (*http.Response, error) {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, b, jsonData)
			resp, err := httpmock.NewJsonResponse(202, map[string]interface{}{
				"value": "fixed",
			})
			return resp, err
		},
	)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	UpdateInAppNotification(userId, params)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "NotificationAPI request ignored.", "The log message should be %s, got: %v", "NotificationAPI request ignored.", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}

func TestUpdateInAppNotificationFailsWith500(t *testing.T) {
	Init(client_id, client_secret)
	params := InAppNotificationPatchRequest{
		TrackingIds: []string{"track123"},
		Opened:      stringPtr("2024-07-16T08:00:00Z"),
		Clicked:     stringPtr("2024-07-16T08:05:00Z"),
		Archived:    stringPtr("2024-07-16T09:00:00Z"),
		Actioned1:   stringPtr("Action1"),
		Actioned2:   stringPtr("Action2"),
		Reply: &struct {
			Date    string `json:"date"`
			Message string `json:"message"`
		}{
			Date:    "2024-07-16T10:00:00Z",
			Message: "Test reply message",
		},
	}
	userId := "13"
	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PATCH", "https://api.notificationapi.com/client_id/users/"+userId+"/notifications/INAPP_WEB",
		func(req *http.Request) (*http.Response, error) {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, b, jsonData)
			resp, err := httpmock.NewJsonResponse(500, map[string]interface{}{
				"value": "fixed",
			})
			return resp, err
		},
	)
	res := UpdateInAppNotification(userId, params)
	assert.EqualErrorf(t, res, "NotificationAPI request failed.", "The log message should be %s, got: %v", "NotificationAPI request failed.", res)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}

func TestUpdateInAppNotificationPasses(t *testing.T) {
	Init(client_id, client_secret)
	params := InAppNotificationPatchRequest{
		TrackingIds: []string{"track123"},
		Opened:      stringPtr("2024-07-16T08:00:00Z"),
		Clicked:     stringPtr("2024-07-16T08:05:00Z"),
		Archived:    stringPtr("2024-07-16T09:00:00Z"),
		Actioned1:   stringPtr("Action1"),
		Actioned2:   stringPtr("Action2"),
		Reply: &struct {
			Date    string `json:"date"`
			Message string `json:"message"`
		}{
			Date:    "2024-07-16T10:00:00Z",
			Message: "Test reply message",
		},
	}
	userId := "13"
	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PATCH", "https://api.notificationapi.com/client_id/users/"+userId+"/notifications/INAPP_WEB",
		func(req *http.Request) (*http.Response, error) {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, b, jsonData)
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"value": "fixed",
			})
			return resp, err
		},
	)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	UpdateInAppNotification(userId, params)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "", "The log message should be %s, got: %v", "", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}

// Helper function to get a pointer to a string
func stringPtr(s string) *string {
	return &s
}
