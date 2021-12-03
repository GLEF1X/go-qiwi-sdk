package types

import (
	"time"
)

type BillStatusValue string

const (
	StatusWaiting  BillStatusValue = "WAITING"
	StatusPaid     BillStatusValue = "PAID"
	StatusRejected BillStatusValue = "REJECTED"
	StatusExpired  BillStatusValue = "EXPIRED"
)

type Customer struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Account string `json:"account"`
}

type CustomFields struct {
	PaySourceFilter string `json:"paySourceFilter"`
	ThemeCode       string `json:"themeCode"`
}

type BillStatus struct {
	Value           BillStatusValue `json:"value"`
	ChangedDateTime time.Time       `json:"changedDateTime"`
}

type Bill struct {
	Amount       RequestAmount `json:"amount"`
	Status       BillStatus    `json:"status"`
	Customer     Customer      `json:"customer"`
	CustomFields CustomFields  `json:"customFields"`
	SiteId       string        `json:"siteId"`
	BillId       string        `json:"billId"`
	PayUrl       string        `json:"payUrl"`
	Comment      string        `json:"comment"`
	CreatedAt    time.Time     `json:"creationDateTime"`
	ExpireAt     time.Time     `json:"expirationDateTime"`
}
