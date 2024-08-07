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

func TestQueryLogsPassesWith202(t *testing.T) {
	Init(client_id, client_secret)
	params := QueryLogsRequest{
		DateRangeFilter: &DateRangeFilter{
			StartTime: 1719600830559, // Example timestamp
			EndTime:   1719600840559, // Example timestamp
		},
		NotificationFilter: []string{"order_tracking"},
		ChannelFilter:      []string{"EMAIL"},
		UserFilter:         []string{"abcd-1234"},
		StatusFilter:       []string{"SUCCESS"},
		TrackingIds:        []string{"172cf2f4-18cd-4f1f-b2ac-e50c7d71891c"},
		RequestFilter:      []string{`request.mergeTags.item="Krabby Patty Burger"`},
		EnvIdFilter:        []string{"6ok6imq9unr2budgiebjdaa6oi"},
	}

	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.notificationapi.com/client_id/logs/query",
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
	QueryLogs(params)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "NotificationAPI request ignored.", "The log message should be %s, got: %v", "NotificationAPI request ignored.", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}

func TestQueryLogsFailsWith500(t *testing.T) {
	Init(client_id, client_secret)
	params := QueryLogsRequest{
		DateRangeFilter: &DateRangeFilter{
			StartTime: 1719600830559, // Example timestamp
			EndTime:   1719600840559, // Example timestamp
		},
		NotificationFilter: []string{"order_tracking"},
		ChannelFilter:      []string{"EMAIL"},
		UserFilter:         []string{"abcd-1234"},
		StatusFilter:       []string{"SUCCESS"},
		TrackingIds:        []string{"172cf2f4-18cd-4f1f-b2ac-e50c7d71891c"},
		RequestFilter:      []string{`request.mergeTags.item="Krabby Patty Burger"`},
		EnvIdFilter:        []string{"6ok6imq9unr2budgiebjdaa6oi"},
	}

	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.notificationapi.com/client_id/logs/query",
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
	res := QueryLogs(params)
	assert.EqualErrorf(t, res, "NotificationAPI request failed.", "The log message should be %s, got: %v", "NotificationAPI request failed.", res)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}

func TestQueryLogsPasses(t *testing.T) {
	Init(client_id, client_secret)
	params := QueryLogsRequest{
		DateRangeFilter: &DateRangeFilter{
			StartTime: 1719600830559, // Example timestamp
			EndTime:   1719600840559, // Example timestamp
		},
		NotificationFilter: []string{"order_tracking"},
		ChannelFilter:      []string{"EMAIL"},
		UserFilter:         []string{"abcd-1234"},
		StatusFilter:       []string{"SUCCESS"},
		TrackingIds:        []string{"172cf2f4-18cd-4f1f-b2ac-e50c7d71891c"},
		RequestFilter:      []string{`request.mergeTags.item="Krabby Patty Burger"`},
		EnvIdFilter:        []string{"6ok6imq9unr2budgiebjdaa6oi"},
	}

	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.notificationapi.com/client_id/logs/query",
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
	QueryLogs(params)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "", "The log message should be %s, got: %v", "", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}
