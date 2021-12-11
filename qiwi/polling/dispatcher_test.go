package polling

import (
	"os"
	"testing"
	"time"

	"github.com/GLEF1X/go-qiwi-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/GLEF1X/go-qiwi-sdk/qiwi"
)

type Setup struct {
	dispatcher *Dispatcher
	apiClient  *qiwi.APIClient
}

func TestNewDispatcherDefaults(t *testing.T) {
	dp := NewDispatcher(&Config{})

	assert.Equal(t, dp.config.MaxConcurrent, defaultConcurrentMaximum)
	assert.Equal(t, dp.config.PollingTimeout, defaultPollingTimeout)
}

func TestDispatcherStartPolling(t *testing.T) {
	s := setup(t)
	go s.dispatcher.StartPolling(s.apiClient)
	s.dispatcher.StopPolling()
}

func TestDispatcherProcessUpdatesProperly(t *testing.T) {
	s := setup(t)

	transactionHandler := func(txn *types.Transaction) {
		assert.Equal(t, 1, txn.ID)
	}
	s.dispatcher.HandleTransaction(transactionHandler)

	s.dispatcher.PropagateEventToHandlers(&types.Transaction{ID: 1})
}

func TestDispatcherErrorHandling(t *testing.T) {
	s := setup(t)

	txnHandler := func(txn *types.Transaction) {
		panic("error")
	}
	s.dispatcher.HandleTransaction(txnHandler)

	s.dispatcher.HandleError(func(e error) {
		assert.Equal(t, "error", e.Error())
	})

	s.dispatcher.PropagateEventToHandlers(&types.Transaction{ID: 1})
}

func TestHandlerProcessEvent(t *testing.T) {
	handler := Handler{fn: func(event interface{}) {
		t.Fail()
	}}

	ch := make(chan struct{})

	go func() {
		time.Sleep(5 * time.Millisecond)
		handler.ProcessEvent(struct{}{}, ch)
	}()
	ch <- struct{}{}
}

func setup(t *testing.T) *Setup {
	t.Helper()
	setup := &Setup{}

	// On CI we substitute values to env
	token, phoneNumber := os.Getenv("API_ACCESS_TOKEN"), os.Getenv("PHONE_NUMBER")
	require.NotEmptyf(t, token, "APIClient token was not set in env")
	require.NotEmptyf(t, phoneNumber, "Phone number was not set in env")

	config, err := qiwi.NewConfig(token, phoneNumber)
	require.NoError(t, err)

	setup.apiClient = qiwi.NewAPIClient(config)
	setup.dispatcher = NewDispatcher(&Config{})

	t.Cleanup(func() {
		setup.apiClient.Close()
	})
	return setup
}
