package NotificationAPI

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSchedulePassesWith202(t *testing.T) {

	Init(client_id, client_secret)
	TrackingId := "TrackingId"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://api.notificationapi.com/client_id/schedule/"+TrackingId,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(202, map[string]interface{}{
				"value": "fixed",
			})
			return resp, err
		},
	)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	DeleteSchedule(TrackingId)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "NotificationAPI request ignored.", "The log message should be %s, got: %v", "NotificationAPI request ignored.", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}
func TestDeleteScheduleFailsWith500(t *testing.T) {
	Init(client_id, client_secret)
	TrackingId := "TrackingId"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://api.notificationapi.com/client_id/schedule/"+TrackingId,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(500, map[string]interface{}{
				"value": "fixed",
			})
			return resp, err
		},
	)
	res := DeleteSchedule(TrackingId)
	assert.EqualErrorf(t, res, "NotificationAPI request failed.", "The log message should be %s, got: %v", "NotificationAPI request failed.", res)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}
func TestDeleteSchedulePasses(t *testing.T) {

	Init(client_id, client_secret)
	TrackingId := "TrackingId"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://api.notificationapi.com/client_id/schedule/"+TrackingId,
		func(req *http.Request) (*http.Response, error) {

			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"value": "fixed",
			})
			return resp, err
		},
	)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	DeleteSchedule(TrackingId)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "", "The log message should be %s, got: %v", "", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}
