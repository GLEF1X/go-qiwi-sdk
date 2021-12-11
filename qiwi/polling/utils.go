package polling

import (
	"sort"

	"github.com/GLEF1X/go-qiwi-sdk/types"
)

func sortHistoryByIDAscending(history *types.History) {
	sort.Slice(history.Transactions, func(i, j int) bool {
		return history.Transactions[i].ID < history.Transactions[j].ID
	})
}
