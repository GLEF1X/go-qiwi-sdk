package p2p

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/GLEF1X/go-qiwi-sdk/types"

	"github.com/stretchr/testify/require"
)

type Setup struct {
	*APIClient
}

func TestAPIClient_CreateBill(t *testing.T) {
	s := setup(t)

	opts := &BillCreationOptions{Amount: types.RequestAmount{Value: 5, Currency: "RUB"}}

	bill, err := s.CreateBill(context.Background(), opts)

	assert.NoError(t, err)
	assert.IsType(t, &types.Bill{}, bill)
}

func TestAPIClient_RejectBill(t *testing.T) {
	s := setup(t)

	opts := &BillCreationOptions{Amount: types.RequestAmount{Value: 5, Currency: "RUB"}}

	bill, err := s.CreateBill(context.Background(), opts)

	assert.NoError(t, err)
	assert.IsType(t, &types.Bill{}, bill)

	err = s.RejectBill(context.Background(), bill.ID)
	assert.NoError(t, err)

	status, err := s.GetBillStatus(context.Background(), bill.ID)
	assert.NoError(t, err)

	assert.Equal(t, types.StatusRejected, status)
}

func TestAPIClient_GetBillStatus(t *testing.T) {
	s := setup(t)

	opts := &BillCreationOptions{Amount: types.RequestAmount{Value: 5, Currency: "RUB"}}

	bill, err := s.CreateBill(context.Background(), opts)
	assert.NoError(t, err)
	assert.IsType(t, &types.Bill{}, bill)

	billStatus, err := s.GetBillStatus(context.Background(), bill.ID)

	assert.NoError(t, err)
	assert.Equal(t, types.StatusWaiting, billStatus)
}

func setup(t *testing.T) *Setup {
	t.Helper()
	setup := &Setup{}

	secretP2PToken := os.Getenv("SECRET_P2P")
	require.NotEmptyf(t, secretP2PToken, "Secret p2p token was not set in env")

	cfg, err := NewConfig(secretP2PToken)
	require.NoErrorf(t, err, "Invalid format of token")
	require.IsType(t, &Config{}, cfg)

	setup.APIClient = NewAPIClient(cfg)

	t.Cleanup(func() {
		setup.APIClient.Close()
	})

	return setup
}
