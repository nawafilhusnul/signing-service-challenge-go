package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

type Repository interface {
	Create(device *domain.Device) error
	Get(id string) (*domain.Device, error)
	FindAll() ([]*domain.Device, error)
}
