package api

import (
	"context"
	"github.com/GLEF1X/qiwi-golang-sdk/types"
	"time"
)

type QiwiPoller struct {
	Timeout      time.Duration
	LastUpdateID int
}

func (p *QiwiPoller) Poll(client *QiwiClient, dest chan types.Transaction, stop chan struct{}) {
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
