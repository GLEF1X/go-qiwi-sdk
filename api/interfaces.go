package api

import "github.com/GLEF1X/go-qiwi-sdk/types"

type Poller interface {
	Poll(client *Client, dest chan types.Transaction, stop chan struct{})
}
