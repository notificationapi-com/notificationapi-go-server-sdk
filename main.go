package NotificationAPI

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

var __client_id, __client_secret string

type EmailAttachments struct {
	Filename string `json:"filename,omitempty"`
	Url      string `json:"url,omitempty"`
}
type SendRequestEmailOptions struct {
	ReplyToAddresses []string           `json:"replyToAddresses,omitempty"`
	CcAddresses      []string           `json:"ccAddresses,omitempty"`
	BccAddresses     []string           `json:"bccAddresses,omitempty"`
	FromName		 string				`json:"fromName,omitempty"`
	FromAddress		 string				`json:"fromAddress,omitempty"`
	Attachments      []EmailAttachments `json:"attachments,omitempty"`
}
type SendRequestApnOptions struct {
	Expiry           *string `json:"expiry,omitempty"`
	Priority         *int    `json:"priority,omitempty"`
	CollapseId       *string `json:"collapseId,omitempty"`
	ThreadId         *string `json:"threadId,omitempty"`
	Badge            *int    `json:"badge,omitempty"`
	Sound            *string `json:"sound,omitempty"`
	ContentAvailable *bool   `json:"contentAvailable,omitempty"`
}
type FcmAndroidOptions struct {
	CollapseKey           *string `json:"collapseKey,omitempty"`
	Priority              *string `json:"priority,omitempty"`
	Ttl                   *int    `json:"ttl,omitempty"`
	RestrictedPackageName *string `json:"restrictedPackageName,omitempty"`
}
type SendRequestFcmOptions struct {
	Android *FcmAndroidOptions `json:"android,omitempty"`
}
type SendRequestOptions struct {
	Email *SendRequestEmailOptions `json:"email,omitempty"`
	Apn   *SendRequestApnOptions   `json:"apn,omitempty"`
	Fcm   *SendRequestFcmOptions   `json:"fcm,omitempty"`
}
type UserPushTokenDevice struct {
	App_id       *string `json:"app_id,omitempty"`
	Ad_id        *string `json:"ad_id,omitempty"`
	Device_id    *string `json:"device_id,omitempty"`
	Platform     *string `json:"platform,omitempty"`
	Manufacturer *string `json:"manufacturer,omitempty"`
	Model        *string `json:"model,omitempty"`
}
type UserPushToken struct {
	Type   *string              `json:"type,omitempty"`
	Token  *string              `json:"token,omitempty"`
	Device *UserPushTokenDevice `json:"device,omitempty"`
}
type UserWebPushToken struct {
	Sub struct {
		Endpoint string `json:"endpoint,omitempty"`
		Keys     struct {
			P256DH string `json:"p256dh,omitempty"`
			Auth   string `json:"auth,omitempty"`
		} `json:"keys,omitempty"`
	} `json:"sub,omitempty"`
}

type User struct {
	Id            string              `json:"id,omitempty"`
	Email         string              `json:"email,omitempty"`
	Number        *string             `json:"number,omitempty"`
	PushTokens    *[]UserPushToken    `json:"pushTokens,omitempty"`
	WebPushTokens *[]UserWebPushToken `json:"webPushTokens,omitempty"` // Added WebPushTokens
}
type SendRequest struct {
	NotificationId    string                 `json:"notificationId"`
	User              User                   `json:"user"`
	MergeTags         map[string]interface{} `json:"mergeTags,omitempty"`
	Replace           map[string]string      `json:"replace,omitempty"`
	ForceChannels     []string               `json:"forceChannels,omitempty"`
	Schedule          string                 `json:"schedule,omitempty"`
	TemplateID        *string                `json:"templateId,omitempty"`
	SubNotificationId *string                `json:"subNotificationId,omitempty"`
	Options           *SendRequestOptions    `json:"options,omitempty"`
}
type UpdateScheduleRequest struct {
	NotificationId    string                 `json:"notificationId,omitempty"`
	User              *User                   `json:"user,omitempty"`
	MergeTags         map[string]interface{} `json:"mergeTags,omitempty"`
	Replace           map[string]string      `json:"replace,omitempty"`
	ForceChannels     []string               `json:"forceChannels,omitempty"`
	Schedule          string                 `json:"schedule,omitempty"`
	TemplateID        *string                `json:"templateId,omitempty"`
	SubNotificationId *string                `json:"subNotificationId,omitempty"`
	Options           *SendRequestOptions    `json:"options,omitempty"`
}
type RetractRequest struct {
	NotificationId string `json:"notificationId,omitempty"`
	UserId         string `json:"userId,omitempty"`
}
type CreateSubNotificationRequest struct {
	NotificationId    string `json:"userId,omitempty"`
	Title             string `json:"title,omitempty"`
	SubNotificationId string `json:"subNotificationId,omitempty"`
}
type DeleteSubNotificationRequest struct {
	NotificationId    string `json:"notificationId,omitempty"`
	SubNotificationId string `json:"subNotificationId,omitempty"`
}
type SetUserPreferencesRequest struct {
	NotificationId    string `json:"notificationId,omitempty"`
	Channel           string `json:"channel,omitempty"`
	State             bool   `json:"state,omitempty"`
	SubNotificationId string `json:"subNotificationId,omitempty"`
}
type UserData struct {
	Email      *string          `json:"email,omitempty"`
	Number     *string          `json:"number,omitempty"`
	PushTokens *[]UserPushToken `json:"pushTokens,omitempty"`
}

func Init(client_id, client_secret string) error {
	if client_id == "" {
		return errors.New("Bad client_id")
	}
	if client_secret == "" {
		return errors.New("Bad client_secret")
	}
	__client_id = client_id
	__client_secret = client_secret
	return nil

}
func basicAuth(client_id, client_secret string) string {
	auth := client_id + ":" + client_secret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}
func request(client *http.Client, method, uri string, data *bytes.Buffer, customAuthorization ...string) error {
	endpoint := "https://api.notificationapi.com/" + __client_id + "/" + uri
	req, err := http.NewRequest(method, endpoint, data)
	if err != nil {
		return fmt.Errorf("error occurred while creating request: %w", err)
	}

	authHeader := "Basic " + basicAuth(__client_id, __client_secret)
	if len(customAuthorization) > 0 {
		authHeader = customAuthorization[0]
	}
	req.Header.Add("Authorization", authHeader)
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}
	if response.StatusCode == 202 {
		fmt.Printf("NotificationAPI request ignored.")
	}

	if response.StatusCode == 500 {
		return errors.New("NotificationAPI request failed.")
	}

	return nil

}

func Send(params SendRequest) error {
	c := httpClient()
	sendRequest, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}
	return request(c, http.MethodPost, "sender", bytes.NewBuffer(sendRequest))
}
func Retract(params RetractRequest) error {
	c := httpClient()
	retractRequest, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}
	return request(c, http.MethodPost, "sender/retract", bytes.NewBuffer(retractRequest))
}
func CreateSubNotification(params CreateSubNotificationRequest) error {
	c := httpClient()
	createSubNotificationRequest, err := json.Marshal(map[string]string{"title": params.Title})
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}
	return request(c, http.MethodPut, "notifications/"+params.NotificationId+"/subNotifications/"+params.SubNotificationId, bytes.NewBuffer(createSubNotificationRequest))
}
func UpdateSchedule(TrackingId string, UpdateScheduleRequest UpdateScheduleRequest) error {
	c := httpClient()
	updateScheduleRequest, err := json.Marshal(UpdateScheduleRequest)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}
	return request(c, http.MethodPatch, "schedule/"+TrackingId, bytes.NewBuffer(updateScheduleRequest))
}
func DeleteSchedule(TrackingId string) error {
	c := httpClient()
	return request(c, http.MethodDelete, "schedule/"+TrackingId, bytes.NewBuffer(nil))
}
func DeleteSubNotification(params DeleteSubNotificationRequest) error {
	c := httpClient()
	return request(c, http.MethodDelete, "notifications/"+params.NotificationId+"/subNotifications/"+params.SubNotificationId, bytes.NewBuffer(nil))
}
func SetUserPreferences(userId string, params []SetUserPreferencesRequest) error {
	c := httpClient()
	setUserPreferencesRequest, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}
	return request(c, http.MethodPost, "user_preferences/"+userId, bytes.NewBuffer(setUserPreferencesRequest))
}

// IdentifyUser hashes the user's ID and sends a POST request with user data
func IdentifyUser(user User) error {
	// Hash the user's ID
	h := hmac.New(sha256.New, []byte(__client_secret))
	h.Write([]byte(user.Id))
	hashedUserID := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Construct custom authorization header
	customAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s:%s", __client_id, user.Id, hashedUserID)))

	userData := UserData{
		Email:      &user.Email,
		Number:     user.Number,
		PushTokens: user.PushTokens,
	}

	data, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("error marshalling user data: %w", err)
	}
	// Use the updated request function with custom authorization
	client := httpClient()
	return request(client, "POST", "users/"+url.QueryEscape(user.Id), bytes.NewBuffer(data), customAuth)
}
