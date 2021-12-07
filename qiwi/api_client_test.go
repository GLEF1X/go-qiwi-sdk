package qiwi

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/GLEF1X/go-qiwi-sdk/qiwi/filters"

	"github.com/GLEF1X/go-qiwi-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

type Setup struct {
	*APIClient
}

func TestAPIClient_RetrieveHistory(t *testing.T) {
	currentTime := time.Now()
	truncatedFor15Minutes := currentTime.Truncate(15 * time.Minute)

	testCases := map[string]*filters.HistoryFilter{
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

		history, err := s.RetrieveHistory(context.Background(), testFilter)
		assert.NoError(t, err)
		assert.IsType(t, &types.History{}, history)
	}
}

func TestAPIClient_GetProfile(t *testing.T) {
	s := setup(t)

	profile, err := s.GetProfile(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &types.Profile{}, profile)
}

func setup(t *testing.T) *Setup {
	t.Helper()
	setup := &Setup{}

	// On CI we substitute values to env
	token, phoneNumber := os.Getenv("API_ACCESS_TOKEN"), os.Getenv("PHONE_NUMBER")
	require.NotEmptyf(t, token, "APIClient token was not set in env")
	require.NotEmptyf(t, phoneNumber, "Phone number was not set in env")

	config, err := NewConfig(token, phoneNumber)
	require.NoError(t, err)

	setup.APIClient = NewAPIClient(config)

	t.Cleanup(func() {
		setup.APIClient.Close()
	})
	return setup
}
