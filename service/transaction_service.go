package service

import (
	"encoding/base64"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

type TransactionService interface {
	SignTransaction(deviceID string, data string) (*domain.SignatureResult, error)
}

type transactionService struct {
	repository persistence.Repository
}

func NewTransactionService(repository persistence.Repository) TransactionService {
	return &transactionService{repository: repository}
}

func (s *transactionService) SignTransaction(deviceID string, data string) (*domain.SignatureResult, error) {
	var result *domain.SignatureResult

	_, err := s.repository.UpdateAtomic(deviceID, func(device *domain.Device) error {
		signer, err := crypto.NewSignerFromDevice(device.Algorithm, []byte(device.PrivateKey))
		if err != nil {
			return err
		}
		securedData := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, data, device.LastSignature)

		signBytes, err := signer.Sign([]byte(securedData))
		if err != nil {
			return err
		}
		signatureBase64 := base64.RawStdEncoding.EncodeToString(signBytes)

		device.SignatureCounter++
		device.LastSignature = signatureBase64

		result = &domain.SignatureResult{
			Signature:  signatureBase64,
			SignedData: securedData,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
