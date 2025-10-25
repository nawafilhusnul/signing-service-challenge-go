package domain

type Device struct {
	ID string `json:"id"`
}

type SignatureResult struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}
