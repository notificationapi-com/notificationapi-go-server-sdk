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



func TestCreateSubNotificationPassesWith202(t *testing.T) {
    Init(client_id,client_secret)
     params:= CreateSubNotificationRequest{NotificationId: "baaz",Title:"test",SubNotificationId:"123"}
    jsonData, _ := json.Marshal(map[string]string{ "title": params.Title })
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder("PUT", "https://api.notificationapi.com/client_id/"+"notifications/"+params.NotificationId+"/subNotifications/"+params.SubNotificationId,
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
    CreateSubNotification(params)
    w.Close()
    out, _ := ioutil.ReadAll(r)
    os.Stdout = rescueStdout
    assert.Equal(t, string(out), "NotificationAPI request ignored.", "The log message should be %s, got: %v", "NotificationAPI request ignored.", out)
    assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
    httpmock.Deactivate()
}
func TestCreateSubNotificationFailsWith500(t *testing.T) {
    Init(client_id,client_secret)
    params:= CreateSubNotificationRequest{NotificationId: "baaz",Title:"test",SubNotificationId:"123"}
    jsonData, _ := json.Marshal(map[string]string{ "title": params.Title })
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder("PUT", "https://api.notificationapi.com/client_id/"+"notifications/"+params.NotificationId+"/subNotifications/"+params.SubNotificationId,
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
    res:=CreateSubNotification(params)
    assert.EqualErrorf(t,    res , "NotificationAPI request failed.", "The log message should be %s, got: %v", "NotificationAPI request failed.",res)
    assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
    httpmock.Deactivate()
}
func TestCreateSubNotificationPasses(t *testing.T) {
    Init(client_id,client_secret)
    params:= CreateSubNotificationRequest{NotificationId: "baaz",Title:"test",SubNotificationId:"123"}
    jsonData, _ := json.Marshal(map[string]string{ "title": params.Title })
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder("PUT", "https://api.notificationapi.com/client_id/"+"notifications/"+params.NotificationId+"/subNotifications/"+params.SubNotificationId,
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
    CreateSubNotification(params)
    w.Close()
    out, _ := ioutil.ReadAll(r)
    os.Stdout = rescueStdout
    assert.Equal(t, string(out), "", "The log message should be %s, got: %v", "", out)
    assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Error should be: %v, got: %v", 1, httpmock.GetTotalCallCount())
    httpmock.Deactivate()
}