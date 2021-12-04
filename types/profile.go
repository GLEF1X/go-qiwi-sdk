package types

import (
	"time"
)

type Profile struct {
	Auth     AuthInfo     `json:"authInfo,omitempty"`
	Contract ContractInfo `json:"contractInfo,omitempty"`
	User     UserInfo     `json:"userInfo,omitempty"`
}

type ContractInfo struct {
	IsBlocked          bool                 `json:"blocked"`
	ContractID         int                  `json:"contractId"`
	CreationDate       time.Time            `json:"creationDate"`
	IdentificationInfo []IdentificationInfo `json:"identificationInfo"`
	SMSNotification    SMSNotification      `json:"smsNotification"`
	NickName           NickName             `json:"nickName"`
	Features           []Feature            `json:",omitempty"`
}

type NickName struct {
	Nickname    string `json:"nickname,omitempty"`
	CanChange   bool   `json:"canChange"`
	CanUse      bool   `json:"CanUse"`
	Description string `json:"description,omitempty"`
}

type Feature struct {
	FeatureID    int    `json:"featureId"`
	FeatureValue string `json:"featureValue"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
}

type SMSNotification struct {
	Price     ResponseAmount `json:"price"`
	IsEnabled bool           `json:"enabled"`
	IsActive  bool           `json:"active"`
	EndDate   *time.Time     `json:"endDate,omitempty"`
}

type IdentificationInfo struct {
	BankAlias           string `json:"bankAlias,omitempty"`
	IdentificationLevel string `json:"identificationLevel,omitempty"`
	IsPassportExpired   bool   `json:"passportExpired,omitempty"`
}

type AuthInfo struct {
	Ip            string            `json:"ip"`
	BoundEmail    string            `json:"boundEmail"`
	LastLoginDate time.Time         `json:"lastLoginDate,omitempty"`
	EmailSettings map[string]string `json:"emailSettings,omitempty"`
	MobilePin     MobilePinInfo     `json:"mobilePinInfo"`
	PassInfo      PassInfo          `json:"passInfo"`
	PersonID      int               `json:"personId"`
	PinInfo       struct {
		PinUsed bool `json:"PinUsed"`
	} `json:"pinInfo"`
	RegistrationDate time.Time `json:"registrationDate"`
}

type PassInfo struct {
	LastPassChange string `json:"lastPassChange"`
	NextPassChange string `json:"nextPassChange"`
	IsPasswordUsed bool   `json:"passwordUsed"`
}

type MobilePinInfo struct {
	LastMobilePinChange *time.Time `json:"lastMobilePinChange,omitempty"`
	IsMobilePinUsed     bool       `json:"mobilePinUsed"`
	NextMobilePinChange string     `json:"nextMobilePinChange"`
}

type UserInfo struct {
	DefaultPayCurrency     int    `json:"defaultPayCurrency"`
	DefaultPaySource       int    `json:"defaultPaySource,omitempty"`
	DefaultPayAccountAlias string `json:"defaultPayAccountAlias,omitempty"`
	Email                  string `json:"email,omitempty"`
	FirstTxnID             int    `json:"firstTxnId,omitempty"`
	Language               string `json:"language"`
	PhoneHash              string `json:"phoneHash"`
	IsPromoEnabled         bool   `json:"promoEnabled,omitempty"`
}
