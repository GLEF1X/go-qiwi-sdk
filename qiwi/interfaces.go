package qiwi

import "github.com/GLEF1X/go-qiwi-sdk/types"

type Poller interface {
	Poll(client *APIClient, dest chan types.Transaction, stop chan struct{})
}
