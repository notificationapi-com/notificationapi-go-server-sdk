package NotificationAPI

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)
var client_id = "client_id"
var client_secret = "client_secret"

func TestInitRaisesGivenEmptyClientId(t *testing.T) {
	assert.EqualErrorf(t, Init("",client_secret), "Bad client_id", "Error should be: %v, got: %v", errors.New("Bad client_id"), Init("",client_secret))
	}
func TestInitRaisesGivenEmptySecret(t *testing.T) {
assert.EqualErrorf(t, Init(client_id,""), "Bad client_secret", "Error should be: %v, got: %v", errors.New("Bad client_secret"), Init(client_id,""))
}

func TestInitPassesGivenIdAndSecret(t *testing.T) {
	Init(client_id,client_secret)
}