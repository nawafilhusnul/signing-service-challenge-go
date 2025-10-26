package api

type CreateDeviceRequest struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label,omitempty"`
}

type SignTransactionRequest struct {
	DeviceID string `json:"deviceId"`
	Data     string `json:"data"`
}
