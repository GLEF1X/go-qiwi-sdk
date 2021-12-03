package api

import (
	"context"
	"github.com/GLEF1X/qiwi-golang-sdk/core"
	"github.com/GLEF1X/qiwi-golang-sdk/types"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	baseQIWIUrl = "https://edge.qiwi.com"
)

type QiwiClient struct {
	config   *Config
	client   *core.HttpClient
	poller   Poller
	validate *validator.Validate // cache are saving for multiply calls
}

func NewQiwiClient(config *Config) *QiwiClient {
	return &QiwiClient{
		config:   config,
		client:   core.NewHttpClient(),
		validate: validator.New(),
	}
}

func (api *QiwiClient) BindPoller(p Poller) {
	api.poller = p
}

// History method helps you to receive transactions on the account.
// More detailed documentation: https://developer.qiwi.com/ru/qiwi-wallet-personal/?http#payments_list
func (api *QiwiClient) History(ctx context.Context, historyFilter *HistoryFilter) (*types.History, error) {
	queryParams, err := historyFilter.ConvertToMapWithValidation(api.validate)
	if err != nil {
		return nil, err
	}
	response, err := api.client.SendRequest(
		ctx,
		&core.Request{
			BaseUrl:            baseQIWIUrl,
			APIEndpoint:        endpoints.Transactions,
			HttpMethod:         http.MethodGet,
			AuthorizationToken: api.config.AuthorizationToken,
			Payload: &core.Payload{
				QueryParams:      queryParams,
				UrlConstructArgs: []interface{}{api.config.GetPhoneNumberForAPIRequests()},
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
