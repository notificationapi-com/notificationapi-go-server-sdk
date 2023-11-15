package NotificationAPI

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)


func TestIdentifyUserSuccess(t *testing.T) {
    Init(client_id, client_secret)

    user := User{Id: "testuser", Email: "test@example.com"}
    userData := UserData{
        Email:      &user.Email,
        Number:     user.Number,
        PushTokens: user.PushTokens,
    }

    jsonData, _ := json.Marshal(userData)

    h := hmac.New(sha256.New, []byte(client_secret))
    h.Write([]byte(user.Id))
    hashedUserID := base64.StdEncoding.EncodeToString(h.Sum(nil))
    customAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(client_id+":"+user.Id+":"+hashedUserID))

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder("POST", "https://api.notificationapi.com/client_id/users/"+url.QueryEscape(user.Id),
        func(req *http.Request) (*http.Response, error) {
            // This assertion checks the Authorization header within the scope where 'req' is defined
            assert.Equal(t, req.Header.Get("Authorization"), customAuth, "Authorization header should match")
            b, err := ioutil.ReadAll(req.Body)
            if err != nil {
                panic(err)
            }
            assert.Equal(t, b, jsonData, "Request body should match")
            resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
                "status": "success",
            })
            return resp, err
        },
    )

    err := IdentifyUser(user)
    assert.NoError(t, err, "IdentifyUser should not return an error")
    assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Total call count should be 1")
}


func TestIdentifyUserFailure(t *testing.T) {
    Init(client_id, client_secret)

    user := User{Id: "testuser", Email: "test@example.com"}

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder("POST", "https://api.notificationapi.com/client_id/users/"+url.QueryEscape(user.Id),
        func(req *http.Request) (*http.Response, error) {
            resp, err := httpmock.NewJsonResponse(500, map[string]interface{}{
                "status": "error",
            })
            return resp, err
        },
    )

    err := IdentifyUser(user)
    assert.Error(t, err, "IdentifyUser should return an error for a 500 response")
    assert.Equal(t, httpmock.GetTotalCallCount(), 1, "Total call count should be 1")
}
