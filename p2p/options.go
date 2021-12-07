package p2p

import (
	"errors"
	"fmt"
	"time"

	"github.com/GLEF1X/go-qiwi-sdk/types"
	"github.com/google/uuid"
)

const (
	DefaultExpireDuration = 15 * time.Minute
)

var (
	NormalizationError = errors.New("BillCreationOptions normalize: ")
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
	if opts.Amount.Value <= 0 {
		return nil, fmt.Errorf("%w amount cannot be <=0", NormalizationError)
	}
	if opts.Amount.Currency == "" {
		return nil, fmt.Errorf("%w currency field cannot be empty", NormalizationError)
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
