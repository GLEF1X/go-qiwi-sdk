package api

import (
	"errors"
	"strings"
)

var (
	ErrWrongConfigCredentials = errors.New("API Token is empty string")
)

type Config struct {
	AuthorizationToken string // QIWI API token received from https://qiwi.com/api
	PhoneNumber        string // phone number
}

func (c *Config) GetPhoneNumberForAPIRequests() string {
	return c.PhoneNumber[1:]
}

func NewConfig(APIAccessToken string, PhoneNumber string) (*Config, error) {
	if strings.TrimSpace(APIAccessToken) == "" {
		return nil, ErrWrongConfigCredentials
	}
	if !strings.HasPrefix(PhoneNumber, "+") {
		PhoneNumber = "+" + PhoneNumber
	}
	return &Config{
		AuthorizationToken: APIAccessToken,
		PhoneNumber:        PhoneNumber,
	}, nil
}
