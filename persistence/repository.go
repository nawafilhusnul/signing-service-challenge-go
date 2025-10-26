package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

type Repository interface {
	Create(device *domain.Device) error
	GetByID(id string) (*domain.Device, error)
	FindAll() ([]*domain.Device, error)
	Update(device *domain.Device) error
	UpdateAtomic(deviceID string, updateFn func(*domain.Device) error) (*domain.Device, error)
}
