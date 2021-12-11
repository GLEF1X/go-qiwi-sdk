package p2p

import (
	"encoding/base64"
	"fmt"

	"github.com/goccy/go-json"
)

type tokenPayload struct {
	Version string
	Data    struct {
		PayinMerchantSiteUID string
		UserID               string
		Secret               string
	}
}

type Config struct {
	SecretToken string // QIWI P2P secret token received from https://p2p.qiwi.com/
}

func NewConfig(secretToken string) (*Config, error) {
	if !isTokenValid(secretToken) {
		return nil, fmt.Errorf("p2p token is invalid")
	}
	return &Config{SecretToken: secretToken}, nil
}

func isTokenValid(plainStringToken string) bool {
	decodedToken, err := base64.StdEncoding.DecodeString(plainStringToken)
	if err != nil {
		return false
	}
	var payload tokenPayload
	if err := json.Unmarshal(decodedToken, &payload); err != nil {
		return false
	}
	return true
}
