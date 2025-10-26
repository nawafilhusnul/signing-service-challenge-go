package domain

import "time"

type Device struct {
	ID               string    `json:"id"`
	Algorithm        string    `json:"algorithm"`
	Label            string    `json:"label"`
	SignatureCounter int       `json:"signatureCounter"`
	LastSignature    string    `json:"-"`
	PrivateKey       string    `json:"-"`
	PublicKey        string    `json:"publicKey"`
	CreatedAt        time.Time `json:"createdAt"`
}

type SignatureResult struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signedData"`
}
