package api

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/GLEF1X/go-qiwi-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

type Setup struct {
	*Client
}

func TestQiwiClient_History(t *testing.T) {
	currentTime := time.Now()
	truncatedFor15Minutes := currentTime.Truncate(15 * time.Minute)

	testCases := map[string]*HistoryFilter{
		"without any filters": {},
		"with rows filter":    {Rows: 50},
		"with start date and end date filter": {
			StartDate: &truncatedFor15Minutes,
			EndDate:   &currentTime,
		},
	}

	s := setup(t)

	for testCaseName, testFilter := range testCases {
		t.Logf("Running test case %s", testCaseName)

		history, err := s.History(context.Background(), testFilter)
		require.NoError(t, err)
		assert.IsType(t, &types.History{}, history)
	}
}

func TestQiwiClient_BindPoller(t *testing.T) {
	s := setup(t)

	qiwiPoller := &QiwiPoller{}
	s.BindPoller(qiwiPoller)

	assert.Equal(t, qiwiPoller, s.poller)
}

func setup(t *testing.T) *Setup {
	t.Helper()
	setup := &Setup{}

	// On CI we substitute values to env
	token, phoneNumber := os.Getenv("API_ACCESS_TOKEN"), os.Getenv("PHONE_NUMBER")
	require.NotEmptyf(t, token, "API token was not set in env")
	require.NotEmptyf(t, phoneNumber, "Phone number was not set in env")

	config, err := NewConfig(token, phoneNumber)
	require.NoError(t, err)

	setup.Client = NewClient(config)

	t.Cleanup(func() {
		setup.Client.Close()
	})
	return setup
}
