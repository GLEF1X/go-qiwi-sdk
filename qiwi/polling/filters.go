package polling

import "github.com/GLEF1X/go-qiwi-sdk/types"

type Filter interface {
	Check(event interface{}) bool
}

// Built-in filters
type transactionFilter struct{}

func (txnFilter transactionFilter) Check(event interface{}) bool {
	_, ok := event.(*types.Transaction)
	return ok
}

type errorFilter struct{}

func (f errorFilter) Check(event interface{}) bool {
	_, ok := event.(error)
	return ok
}
