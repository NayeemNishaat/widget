package lib

type StripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type Response struct {
	Error   bool           `json:"error"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}
