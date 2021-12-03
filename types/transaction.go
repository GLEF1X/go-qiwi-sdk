package types

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID         int            `json:"txnId"`
	PersonId   int            `json:"personId"`
	Date       time.Time      `json:"date"`
	Status     string         `json:"status"`
	TxnType    string         `json:"type"`
	Amount     ResponseAmount `json:"sum"`
	Comment    string         `json:"comment"`
	Commission ResponseAmount `json:"commission"`
}

func (t Transaction) String() string {
	return fmt.Sprintf("№%d - %s - %f - %s", t.ID, t.Status, t.Amount.Value, t.Comment)
}

type History struct {
	Transactions []Transaction `json:"data"`
	NextTxnId    int           `json:"nextTxnId"`
	NextTxnDate  time.Time     `json:"nextTxnDate"`
}
