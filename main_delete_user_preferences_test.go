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

func TestDeleteUserPreferencesPassesWith202(t *testing.T) {
	Init(client_id, client_secret)
	userId := "123"
	notificationId := "baaz"
	subNotificationId := "123"
	params := map[string]string{"notificationId": notificationId, "subNotificationId": subNotificationId}
	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://api.notificationapi.com/client_id/users/"+userId+"/preferences",
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
	DeleteUserPreferences(userId, notificationId, subNotificationId)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "NotificationAPI request ignored.", "The log message should be %s, got: %v", "NotificationAPI request ignored.", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}

func TestDeleteUserPreferencesFailsWith500(t *testing.T) {
	Init(client_id, client_secret)
	userId := "13"
	notificationId := "baaz"
	subNotificationId := "123"
	params := map[string]string{"notificationId": notificationId, "subNotificationId": subNotificationId}
	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://api.notificationapi.com/client_id/users/"+userId+"/preferences",
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
	res := DeleteUserPreferences(userId, notificationId, subNotificationId)
	assert.EqualErrorf(t, res, "NotificationAPI request failed.", "The log message should be %s, got: %v", "NotificationAPI request failed.", res)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}

func TestDeleteUserPreferencesPasses(t *testing.T) {
	Init(client_id, client_secret)
	userId := "13"
	notificationId := "baaz"
	subNotificationId := "123"
	params := map[string]string{"notificationId": notificationId, "subNotificationId": subNotificationId}
	jsonData, _ := json.Marshal(params)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://api.notificationapi.com/client_id/users/"+userId+"/preferences",
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
	DeleteUserPreferences(userId, notificationId, subNotificationId)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Equal(t, string(out), "", "The log message should be %s, got: %v", "", out)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
	httpmock.Deactivate()
}