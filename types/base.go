package types

type ResponseAmount struct {
	Value    float64 `json:"amount"`
	Currency int     `json:"currency"`
}

type RequestAmount struct {
	Value    float64 `json:"value,string"`
	Currency string  `json:"currency"`
}
