package p2p

import (
	"context"
	"github.com/GLEF1X/qiwi-golang-sdk/core"
	"github.com/GLEF1X/qiwi-golang-sdk/core/endpoints"
	"github.com/GLEF1X/qiwi-golang-sdk/types"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	baseP2PQiwiAPIUrl = "https://api.qiwi.com"
)

type Client struct {
	config *Config
	client *core.HttpClient
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		client: core.NewHttpClient(),
	}
}

func (api *Client) CreateBill(ctx context.Context, billPayload *BillCreationOptions) (*types.Bill, error) {
	billId := uuid.New().String()
	data, err := api.client.SendRequest(
		ctx,
		&core.Request{
			BaseUrl:            baseP2PQiwiAPIUrl,
			APIEndpoint:        endpoints.CreateBill,
			HttpMethod:         http.MethodPut,
			AuthorizationToken: api.config.SecretToken,
			Payload: &core.Payload{
				UrlConstructArgs: []interface{}{billId},
				Body:             billPayload,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	var bill *types.Bill
	if err := json.Unmarshal(data, &bill); err != nil {
		return nil, err
	}
	return bill, nil
}

func (api *Client) CheckBillStatus(ctx context.Context, billID string) (*types.Bill, bool, error) {
	data, err := api.client.SendRequest(
		ctx,
		&core.Request{
			BaseUrl:            baseP2PQiwiAPIUrl,
			APIEndpoint:        endpoints.CheckBillStatus,
			HttpMethod:         http.MethodGet,
			AuthorizationToken: api.config.SecretToken,
			Payload:            &core.Payload{UrlConstructArgs: []interface{}{billID}},
		},
	)
	if err != nil {
		return nil, false, err
	}
	var bill types.Bill
	if err := json.Unmarshal(data, &bill); err != nil {
		return nil, false, err
	}
	if bill.Status.Value != types.StatusPaid {
		return &bill, false, nil
	}
	return &bill, true, nil
}

func (api *Client) RejectBill(ctx context.Context, billID string) error {
	_, err := api.client.SendRequest(
		ctx,
		&core.Request{
			BaseUrl:            baseP2PQiwiAPIUrl,
			APIEndpoint:        endpoints.RejectBill,
			HttpMethod:         http.MethodPost,
			AuthorizationToken: api.config.SecretToken,
			Payload:            &core.Payload{UrlConstructArgs: []interface{}{billID}},
		},
	)
	if err != nil {
		return err
	}
	return nil
}
