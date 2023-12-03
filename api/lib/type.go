package lib

type StripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type Response struct {
	Error   bool           `json:"error,omitempty"`
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}
