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

type ExpectedUpdateScheduleRequestJsonData struct {
	NotificationId string `json:"notificationId"`
	Schedule       string `json:"schedule"`
}

func TestUpdateSchedulePassesWith202(t *testing.T) {
	Init(client_id, client_secret)
	TrackingId := "TrackingId"
	UpdateScheduleRequest := UpdateScheduleRequest{NotificationId: "baaz", Schedule: "2024-02-20T14:38:03.509Z"}
	jsonData, _ := json.Marshal(ExpectedUpdateScheduleRequestJsonData{NotificationId: "baaz", Schedule: "2024-02-20T14:38:03.509Z"})
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PATCH", "https://api.notificationapi.com/client_id/"+"schedule/"+TrackingId,
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
	UpdateSchedule(TrackingId, UpdateScheduleRequest)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "NotificationAPI request ignored.", "The log message should be %s, got: %v", "NotificationAPI request ignored.", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}
func TestUpdateScheduleFailsWith500(t *testing.T) {
	Init(client_id, client_secret)
	TrackingId := "TrackingId"
	UpdateScheduleRequest := UpdateScheduleRequest{NotificationId: "baaz", Schedule: "2024-02-20T14:38:03.509Z"}
	jsonData, _ := json.Marshal(ExpectedUpdateScheduleRequestJsonData{NotificationId: "baaz", Schedule: "2024-02-20T14:38:03.509Z"})
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PATCH", "https://api.notificationapi.com/client_id/"+"schedule/"+TrackingId,
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
	res := UpdateSchedule(TrackingId, UpdateScheduleRequest)
	assert.EqualErrorf(t, res, "NotificationAPI request failed.", "The log message should be %s, got: %v", "NotificationAPI request failed.", res)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}
func TestUpdateSchedulePasses(t *testing.T) {
	Init(client_id, client_secret)
	TrackingId := "TrackingId"
	UpdateScheduleRequest := UpdateScheduleRequest{NotificationId: "baaz", Schedule: "2024-02-20T14:38:03.509Z"}
	jsonData, _ := json.Marshal(ExpectedUpdateScheduleRequestJsonData{NotificationId: "baaz", Schedule: "2024-02-20T14:38:03.509Z"})
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PATCH", "https://api.notificationapi.com/client_id/"+"schedule/"+TrackingId,
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
	UpdateSchedule(TrackingId, UpdateScheduleRequest)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "", "The log message should be %s, got: %v", "", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}
