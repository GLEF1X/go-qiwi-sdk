package qiwi

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrAPITokenIsEmpty          = errors.New("APIClient Token is empty string")
	ErrPhoneNumberInvalidFormat = errors.New("phone number does not match the required format")
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
		return nil, ErrAPITokenIsEmpty
	}

	if !strings.HasPrefix(PhoneNumber, "+") {
		PhoneNumber = "+" + PhoneNumber
	}

	if !phoneNumberRegexp.MatchString(PhoneNumber) {
		return nil, ErrPhoneNumberInvalidFormat
	}

	return &Config{
		AuthorizationToken: APIAccessToken,
		PhoneNumber:        PhoneNumber,
	}, nil
}
