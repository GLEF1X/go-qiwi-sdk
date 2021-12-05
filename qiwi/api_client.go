package qiwi

import (
	"context"
	"net/http"

	"github.com/GLEF1X/go-qiwi-sdk/qiwi/filters"

	"github.com/GLEF1X/go-qiwi-sdk/core/client"

	"github.com/GLEF1X/go-qiwi-sdk/core/endpoints"
	"github.com/GLEF1X/go-qiwi-sdk/types"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	baseQIWIUrl = "https://edge.qiwi.com"
)

func NewAPIClient(config *Config) *APIClient {
	return &APIClient{
		config: config,
		httpClient: client.NewHttp(
			client.WithBaseURL(baseQIWIUrl),
			client.WithDefaultHeaders(map[string]string{"Authorization": "Bearer " + config.AuthorizationToken}),
		),
		validate: validator.New(),
	}
}

type APIClient struct {
	config     *Config
	httpClient *client.Http
	poller     Poller
	validate   *validator.Validate // internal cache are saving for multiply validations
}

// LoadHistory method helps you to receive transactions on the account.
// More detailed documentation: https://developer.qiwi.com/ru/qiwi-wallet-personal/?http#payments_list
func (c *APIClient) LoadHistory(ctx context.Context, historyFilter *filters.HistoryFilter) (*types.History, error) {
	queryParams, err := historyFilter.ConvertToMapWithValidation(c.validate)
	if err != nil {
		return nil, err
	}
	response, err := c.httpClient.SendRequest(
		ctx,
		&client.Request{
			APIEndpoint: endpoints.GetTransactions,
			HttpMethod:  http.MethodGet,
			Payload: client.Payload{
				QueryParams:      queryParams,
				URLConstructArgs: []interface{}{c.config.GetPhoneNumberForAPIRequests()},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	var history *types.History
	if err := json.Unmarshal(response, &history); err != nil {
		return nil, err
	}
	return history, nil
}

// GetProfile helps you to get profile info. You can find
// detailed documentation here: https://developer.qiwi.com/ru/qiwi-wallet-personal/#profile
func (c *APIClient) GetProfile(ctx context.Context) (*types.Profile, error) {
	response, err := c.httpClient.SendRequest(
		ctx,
		&client.Request{
			APIEndpoint: endpoints.GetProfile,
			HttpMethod:  http.MethodGet,
		},
	)
	if err != nil {
		return nil, err
	}
	var profile *types.Profile
	if err := json.Unmarshal(response, &profile); err != nil {
		return nil, err
	}
	return profile, err
}

func (c *APIClient) BindPoller(p Poller) {
	c.poller = p
}

func (c *APIClient) Close() {
	c.httpClient.Close()
}
