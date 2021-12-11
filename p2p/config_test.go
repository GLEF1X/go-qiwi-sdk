package p2p

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func encodeToBase64(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	defer func() {
		if err := encoder.Close(); err != nil {
			panic(err)
		}
	}()
	err := json.NewEncoder(encoder).Encode(v)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestNewConfig(t *testing.T) {
	mockToken, err := encodeToBase64(tokenPayload{Version: "1", Data: struct {
		PayinMerchantSiteUID string
		UserID               string
		Secret               string
	}{PayinMerchantSiteUID: "test", UserID: "test", Secret: "test"}})
	assert.NoError(t, err)

	_, err = NewConfig(mockToken)
	assert.NoError(t, err)
}

func TestNewConfigFailIfPayloadIsInvalid(t *testing.T) {
	_, err := NewConfig("bla-bla-bla")
	assert.Error(t, err)
	assert.EqualError(t, err, "p2p token is invalid")

}
