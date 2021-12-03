package types

import (
	"time"
)

type BillStatus string

const (
	StatusWaiting  BillStatus = "WAITING"
	StatusPaid     BillStatus = "PAID"
	StatusRejected BillStatus = "REJECTED"
	StatusExpired  BillStatus = "EXPIRED"
)

type Bill struct {
	Amount       RequestAmount `json:"amount"`
	Status       Status        `json:"status"`
	Customer     Customer      `json:"customer"`
	CustomFields CustomFields  `json:"customFields"`
	CreatedAt    time.Time     `json:"creationDateTime"`
	ExpireAt     time.Time     `json:"expirationDateTime"`
	SiteId       string        `json:"siteId"`
	ID           string        `json:"billId"`
	PayUrl       string        `json:"payUrl"`
	Comment      string        `json:"comment"`
}

type Customer struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Account string `json:"account"`
}

type CustomFields struct {
	PaySourceFilter string `json:"paySourceFilter"`
	ThemeCode       string `json:"themeCode"`
}

type Status struct {
	Value           BillStatus `json:"value"`
	ChangedDateTime time.Time  `json:"changedDateTime"`
}
