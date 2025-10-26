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
	device, err := s.repository.GetByID(deviceID)
	if err != nil {
		return nil, err
	}

	signer, err := crypto.NewSignerFromDevice(device.Algorithm, []byte(device.PrivateKey))
	if err != nil {
		return nil, err
	}

	securedData := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, data, device.LastSignature)
	signBytes, err := signer.Sign([]byte(securedData))
	if err != nil {
		return nil, err
	}

	signedData := base64.RawStdEncoding.EncodeToString(signBytes)
	device.SignatureCounter++
	device.LastSignature = signedData

	err = s.repository.Update(device)
	if err != nil {
		return nil, err
	}

	return &domain.SignatureResult{
		SignedData: securedData,
	}, nil
}
