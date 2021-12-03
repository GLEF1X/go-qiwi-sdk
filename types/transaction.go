package types

import (
	"fmt"
	"time"
)

type Transaction struct {
	Commission ResponseAmount `json:"commission"`
	Amount     ResponseAmount `json:"sum"`
	Date       time.Time      `json:"date"`
	Status     string         `json:"status"`
	TxnType    string         `json:"type"`
	Comment    string         `json:"comment"`
	ID         int            `json:"txnId"`
	PersonId   int            `json:"personId"`
}

func (t Transaction) String() string {
	return fmt.Sprintf("№%d - %s - %f - %s", t.ID, t.Status, t.Amount.Value, t.Comment)
}

type History struct {
	Transactions []Transaction `json:"data"`
	NextTxnId    int           `json:"nextTxnId"`
	NextTxnDate  time.Time     `json:"nextTxnDate"`
}
