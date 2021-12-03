# qiwi-golang-sdk


### Quick start:

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
	client := p2p.NewClient(config)
	bill, err := client.CreateBill(
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
