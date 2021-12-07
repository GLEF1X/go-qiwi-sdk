package p2p

import (
	"context"
	"net/http"

	"github.com/GLEF1X/go-qiwi-sdk/core/client"

	"github.com/GLEF1X/go-qiwi-sdk/core/endpoints"
	"github.com/GLEF1X/go-qiwi-sdk/types"
	"github.com/goccy/go-json"
)

const baseP2PQiwiAPIURL = "https://api.qiwi.com"

type APIClient struct {
	config     *Config
	httpClient *client.Http
}

func NewAPIClient(config *Config) *APIClient {
	return &APIClient{
		config: config,
		httpClient: client.NewHttp(
			client.WithBaseURL(baseP2PQiwiAPIURL),
			client.WithDefaultHeaders(map[string]string{"Authorization": "Bearer " + config.SecretToken}),
		),
	}
}

func (api *APIClient) CreateBill(ctx context.Context, options *BillCreationOptions) (*types.Bill, error) {
	options, err := options.Normalize()
	if err != nil {
		return nil, err
	}
	response, err := api.httpClient.SendRequest(
		ctx,
		&client.Request{
			APIEndpoint: endpoints.CreateBill,
			HttpMethod:  http.MethodPut,
			Payload: client.Payload{
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
		&client.Request{
			APIEndpoint: endpoints.CheckBillStatus,
			HttpMethod:  http.MethodGet,
			Payload:     client.Payload{URLConstructArgs: []interface{}{billID}},
		},
	)
	if err != nil {
		return "", err
	}
	var bill *types.Bill
	if err := json.Unmarshal(response, &bill); err != nil {
		return "", err
	}
	return bill.Status.Value, nil
}

func (api *APIClient) RejectBill(ctx context.Context, billID string) error {
	_, err := api.httpClient.SendRequest(
		ctx,
		&client.Request{
			APIEndpoint: endpoints.RejectBill,
			HttpMethod:  http.MethodPost,
			Payload:     client.Payload{URLConstructArgs: []interface{}{billID}},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (api *APIClient) Close() {
	api.httpClient.Close()
}
