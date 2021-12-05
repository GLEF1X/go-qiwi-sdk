package qiwi

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ConfigInputInvalidErr = errors.New("Config invalid: ")
)

var phoneNumberRegexp *regexp.Regexp

func init() {
	phoneNumberRegexp = regexp.MustCompile(`^[+]?[(]?[0-9]{3}[)]?[-\s.]?[0-9]{3}[-\s.]?[0-9]{4,6}$`)
}

type Config struct {
	AuthorizationToken string // QIWI APIClient token received from https://qiwi.com/api
	PhoneNumber        string // phone number
}

func (c *Config) GetPhoneNumberForAPIRequests() string {
	return c.PhoneNumber[1:]
}

func NewConfig(APIAccessToken string, PhoneNumber string) (*Config, error) {
	if strings.TrimSpace(APIAccessToken) == "" {
		return nil, fmt.Errorf("%w api access token is empty", ConfigInputInvalidErr)
	}

	if !strings.HasPrefix(PhoneNumber, "+") {
		PhoneNumber = "+" + PhoneNumber
	}

	if !phoneNumberRegexp.MatchString(PhoneNumber) {
		return nil, fmt.Errorf("%w wrong phone number format", ConfigInputInvalidErr)
	}

	return &Config{
		AuthorizationToken: APIAccessToken,
		PhoneNumber:        PhoneNumber,
	}, nil
}
