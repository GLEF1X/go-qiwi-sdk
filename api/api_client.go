package api

import (
	"context"
	"net/http"

	"github.com/GLEF1X/qiwi-golang-sdk/core"
	"github.com/GLEF1X/qiwi-golang-sdk/core/endpoints"
	"github.com/GLEF1X/qiwi-golang-sdk/types"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	baseQIWIUrl = "https://edge.qiwi.com"
)

type QiwiClient struct {
	config     *Config
	httpClient *core.HttpClient
	poller     Poller
	validate   *validator.Validate // cache are saving for multiply calls
}

func NewQiwiClient(config *Config) *QiwiClient {
	return &QiwiClient{
		config:     config,
		httpClient: core.NewHttpClient(),
		validate:   validator.New(),
	}
}

func (c *QiwiClient) BindPoller(p Poller) {
	c.poller = p
}

func (c *QiwiClient) Close() {
	c.httpClient.Close()
}

// History method helps you to receive transactions on the account.
// More detailed documentation: https://developer.qiwi.com/ru/qiwi-wallet-personal/?http#payments_list
func (c *QiwiClient) History(ctx context.Context, historyFilter *HistoryFilter) (*types.History, error) {
	queryParams, err := historyFilter.ConvertToMapWithValidation(c.validate)
	if err != nil {
		return nil, err
	}
	response, err := c.httpClient.SendRequest(
		ctx,
		&core.Request{
			BaseUrl:            baseQIWIUrl,
			APIEndpoint:        endpoints.GetTransactions,
			HttpMethod:         http.MethodGet,
			AuthorizationToken: c.config.AuthorizationToken,
			Payload: &core.Payload{
				QueryParams:      queryParams,
				URLConstructArgs: []interface{}{c.config.GetPhoneNumberForAPIRequests()},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	var res *types.History
	if err := json.Unmarshal(response, &res); err != nil {
		return nil, err
	}
	return res, nil
}
