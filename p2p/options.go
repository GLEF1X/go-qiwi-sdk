package p2p

import (
	"github.com/GLEF1X/qiwi-golang-sdk/types"
)

type BillCreationOptions struct {
	Amount  types.RequestAmount `json:"amount"`
	Comment string              `json:"comment,omitempty"`
	// ExpirationDateTime time.Time      `json:"expirationDateTime"`
	Customer types.Customer `json:"customer,omitempty"`
}
