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

func TestSendPassesWith202(t *testing.T) {
	Init(client_id, client_secret)
	params := SendRequest{NotificationId: "baaz"}

	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.notificationapi.com/client_id/sender",
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
	Send(params)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "NotificationAPI request ignored.", "The log message should be %s, got: %v", "NotificationAPI request ignored.", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}
func TestSendFailsWith500(t *testing.T) {
	Init(client_id, client_secret)
	params := SendRequest{NotificationId: "baz"}
	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.notificationapi.com/client_id/sender",
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
	res := Send(params)
	assert.EqualErrorf(t, res, "NotificationAPI request failed.", "The log message should be %s, got: %v", "NotificationAPI request failed.", res)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}
func TestSendPasses(t *testing.T) {
	Init(client_id, client_secret)
	orders := []interface{}{
		map[string]string{"id": "123", "productName": "hasan"},
		map[string]string{"id": "124", "productName": "socks"},
	}

	mergeTags := map[string]interface{}{
		"firstName": "Jane",
		"lastName":  "Doe",
		"orders":    orders,
	}
	params := SendRequest{NotificationId: "baz", User: User{Id: "asd"}, MergeTags: mergeTags}
	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.notificationapi.com/client_id/sender",
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
	Send(params)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "", "The log message should be %s, got: %v", "", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}
