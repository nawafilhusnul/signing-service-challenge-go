package domain

import "time"

type Device struct {
	ID               string    `json:"id"`
	Algorithm        string    `json:"algorithm"`
	Label            string    `json:"label"`
	SignatureCounter int       `json:"signature_counter"`
	LastSignature    string    `json:"-"`
	PrivateKey       []byte    `json:"-"`
	PublicKey        []byte    `json:"public_key"`
	CreatedAt        time.Time `json:"created_at"`
}

type SignatureResult struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}
