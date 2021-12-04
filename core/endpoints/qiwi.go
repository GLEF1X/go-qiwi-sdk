package endpoints

import "fmt"

type Endpoint string

func (e Endpoint) Resolve(args []interface{}) string {
	return fmt.Sprintf(string(e), args...)
}

const (
	// Standard API
	GetTransactions Endpoint = "/payment-history/v2/persons/%s/payments"
	CreateBill      Endpoint = "/partner/bill/v1/bills/%s"
	GetProfile      Endpoint = "/person-profile/v1/profile/current"

	// P2P urls
	CheckBillStatus Endpoint = "/partner/bill/v1/bills/%s"
	RejectBill      Endpoint = "/partner/bill/v1/bills/%s/reject"
)
