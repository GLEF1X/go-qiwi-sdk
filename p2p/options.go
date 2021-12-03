package p2p

import (
	"github.com/GLEF1X/qiwi-golang-sdk/types"
	"github.com/google/uuid"
	"time"
)

const (
	DefaultExpireDuration = 15 * time.Minute
)

type BillCreationOptions struct {
	BillID             string              `json:"-"`
	Amount             types.RequestAmount `json:"amount"`
	Comment            *string             `json:"comment,omitempty"`
	Customer           *types.Customer     `json:"customer,omitempty"`
	ExpirationDateTime *time.Time          `json:"expirationDateTime,omitempty"`
}

func (opts *BillCreationOptions) Normalize() (*BillCreationOptions, error) {
	if opts.BillID == "" {
		opts.BillID = uuid.New().String()
	}
	if opts.ExpirationDateTime == nil {
		if err := opts.SetDefaultExpirationDateTime(); err != nil {
			return nil, err
		}

	}
	return opts, nil
}

func (opts *BillCreationOptions) SetDefaultExpirationDateTime() error {
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}
	t := time.Now().Add(DefaultExpireDuration).In(moscowLocation)
	opts.ExpirationDateTime = &t
	return nil
}
