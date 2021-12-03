package p2p

import (
	"encoding/base64"
	"fmt"
)

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
	newTokenPayload := struct {
		version string
		data    struct {
			payinMerchantSiteUID string
			userID               string
			secret               string
		}
	}{}
	if err := json.Unmarshal(decodedToken, &newTokenPayload); err != nil {
		return false
	}
	return true
}
