package endpoints

import "fmt"

type Endpoint string

func (e Endpoint) Resolve(args []interface{}) string {
	return fmt.Sprintf(string(e), args...)
}

// Standard API
const (
	GetTransactions Endpoint = "/payment-history/v2/persons/%s/payments"
	CreateBill      Endpoint = "/partner/bill/v1/bills/%s"
	GetProfile      Endpoint = "/person-profile/v1/profile/current"
	GetCheque       Endpoint = "/payment-history/v1/transactions/%d/cheque/file"
)

// P2P
const (
	CheckBillStatus Endpoint = "/partner/bill/v1/bills/%s"
	RejectBill      Endpoint = "/partner/bill/v1/bills/%s/reject"
)
