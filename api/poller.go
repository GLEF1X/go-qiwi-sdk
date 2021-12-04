package api

import (
	"context"
	"time"

	"github.com/GLEF1X/go-qiwi-sdk/types"
)

type QiwiPoller struct {
	Timeout      time.Duration
	LastUpdateID int
}

func (p *QiwiPoller) Poll(client *Client, dest chan types.Transaction, stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
		}

		bunchOfTransactions, err := client.History(context.Background(), &HistoryFilter{Rows: 50})
		if err != nil {
			//
		}
		for _, transaction := range bunchOfTransactions.Transactions {
			p.LastUpdateID = transaction.ID
			dest <- transaction
		}

	}
}
