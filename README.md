# SDK for QIWI API, written with golang


### Quick example:

```go
package main

import (
	"context"
	"fmt"
	"github.com/GLEF1X/qiwi-golang-sdk/p2p"
	"github.com/GLEF1X/qiwi-golang-sdk/types"
)

func main() {
	config, err := p2p.NewConfig("some token")
	if err != nil {
		// handle it gracefully
		panic(err)
	}
	httpClient := p2p.NewAPIClient(config)
	bill, err := httpClient.CreateBill(
		context.Background(),
		&p2p.BillCreationOptions{
			Amount: types.RequestAmount{
				Value: 5.0, Currency: "RUB",
			},
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(bill.PayUrl)
}
```

### Polling: 

```go
package main

import (
	"log"

	"github.com/GLEF1X/go-qiwi-sdk/qiwi"
	"github.com/GLEF1X/go-qiwi-sdk/qiwi/polling"
	"github.com/GLEF1X/go-qiwi-sdk/types"
)

func main() {
	config, err := qiwi.NewConfig("token", "+phone_number")
	if err != nil {
		panic(err)
	}

	apiClient := qiwi.NewAPIClient(config)

	dp := polling.NewDispatcher(&polling.Config{MaxConcurrent: 10})

	dp.HandleTransaction(func(transaction *types.Transaction) {
		log.Println(transaction.ID)
	})

	dp.HandleError(func(err error) {
		log.Println("Handle error in error handler")
	})

	dp.StartPolling(apiClient)
}
```
