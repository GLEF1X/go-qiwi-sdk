package qiwi

import (
	"context"
	"time"

	"github.com/GLEF1X/go-qiwi-sdk/qiwi/filters"

	"github.com/GLEF1X/go-qiwi-sdk/types"
)

type QiwiPoller struct {
	Timeout      time.Duration
	LastUpdateID int
}

func (p *QiwiPoller) Poll(client *APIClient, dest chan types.Transaction, stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
		}

		bunchOfTransactions, _ := client.LoadHistory(context.Background(), &filters.HistoryFilter{Rows: 50})
		for _, transaction := range bunchOfTransactions.Transactions {
			p.LastUpdateID = transaction.ID
			dest <- transaction
		}

	}
}
