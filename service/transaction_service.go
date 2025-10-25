package service

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

type TransactionService interface {
	SignTransaction(deviceID string, data []byte) (*domain.SignatureResult, error)
}

type transactionService struct {
	repository persistence.Repository
}

func NewTransactionService(repository persistence.Repository) TransactionService {
	return &transactionService{repository: repository}
}

func (s *transactionService) SignTransaction(deviceID string, data []byte) (*domain.SignatureResult, error) {
	panic("not implemented yet")
}
