package p2p

import (
	"context"
	"net/http"

	"github.com/GLEF1X/go-qiwi-sdk/core"
	"github.com/GLEF1X/go-qiwi-sdk/core/endpoints"
	"github.com/GLEF1X/go-qiwi-sdk/types"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	baseP2PQiwiAPIURL = "https://api.qiwi.com"
)

type APIClient struct {
	config     *Config
	httpClient *core.WrappedHTTPClient
}

func NewAPIClient(config *Config) *APIClient {
	return &APIClient{
		config:     config,
		httpClient: core.NewHttpClient(),
	}
}

func (api *APIClient) Close() {
	api.httpClient.Close()
}

func (api *APIClient) CreateBill(ctx context.Context, options *BillCreationOptions) (*types.Bill, error) {
	options, err := options.Normalize()
	if err != nil {
		return nil, err
	}
	response, err := api.httpClient.SendRequest(
		ctx,
		&core.Request{
			BaseUrl:            baseP2PQiwiAPIURL,
			APIEndpoint:        endpoints.CreateBill,
			HttpMethod:         http.MethodPut,
			AuthorizationToken: api.config.SecretToken,
			Payload: &core.Payload{
				URLConstructArgs: []interface{}{options.BillID},
				Body:             options,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	var bill *types.Bill
	if err := json.Unmarshal(response, &bill); err != nil {
		return nil, err
	}
	return bill, nil
}

func (api *APIClient) GetBillStatus(ctx context.Context, billID string) (types.BillStatus, error) {
	response, err := api.httpClient.SendRequest(
		ctx,
		&core.Request{
			BaseUrl:            baseP2PQiwiAPIURL,
			APIEndpoint:        endpoints.CheckBillStatus,
			HttpMethod:         http.MethodGet,
			AuthorizationToken: api.config.SecretToken,
			Payload:            &core.Payload{URLConstructArgs: []interface{}{billID}},
		},
	)
	if err != nil {
		return "", err
	}
	var bill types.Bill
	if err := json.Unmarshal(response, &bill); err != nil {
		return "", err
	}
	return bill.Status.Value, nil
}

func (api *APIClient) RejectBill(ctx context.Context, billID string) error {
	_, err := api.httpClient.SendRequest(
		ctx,
		&core.Request{
			BaseUrl:            baseP2PQiwiAPIURL,
			APIEndpoint:        endpoints.RejectBill,
			HttpMethod:         http.MethodPost,
			AuthorizationToken: api.config.SecretToken,
			Payload:            &core.Payload{URLConstructArgs: []interface{}{billID}},
		},
	)
	if err != nil {
		return err
	}
	return nil
}
