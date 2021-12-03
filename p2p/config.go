package p2p

import (
	"encoding/base64"
	"fmt"
)

type Config struct {
	SecretToken string // QIWI P2P secret key received from https://p2p.qiwi.com/
}

func NewConfig(secretToken string) (*Config, error) {
	_, err := newDecodedToken(secretToken)
	if err != nil {
		return nil, fmt.Errorf("p2p token is invalid")
	}
	return &Config{SecretToken: secretToken}, nil
}

type decodedP2PToken struct {
	version string
	data    struct {
		payinMerchantSiteUID string
		userID               string
		secret               string
	}
}

func newDecodedToken(plainStringToken string) (*decodedP2PToken, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(plainStringToken)
	if err != nil {
		return nil, err
	}
	var token *decodedP2PToken
	if err := json.Unmarshal(decodedToken, &token); err != nil {
		return nil, err
	}
	return token, nil
}
