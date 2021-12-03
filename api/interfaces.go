package api

import "github.com/GLEF1X/qiwi-golang-sdk/types"

type Poller interface {
	Poll(client *QiwiClient, dest chan types.Transaction, stop chan struct{})
}
