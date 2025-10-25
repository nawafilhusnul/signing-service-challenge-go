package persistence

import (
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type InMemoryRepository struct {
	mu      sync.RWMutex
	devices map[string]*domain.Device
}

func NewInMemoryRepository() Repository {
	return &InMemoryRepository{
		mu:      sync.RWMutex{},
		devices: make(map[string]*domain.Device),
	}
}

func (r *InMemoryRepository) Create(device *domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; exists {
		return domain.ErrDeviceAlreadyExists
	}

	r.devices[device.ID] = device
	return nil
}

func (r *InMemoryRepository) GetByID(id string) (*domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.devices[id]
	if !exists {
		return nil, domain.ErrDeviceNotFound
	}

	return device, nil
}

func (r *InMemoryRepository) FindAll() ([]*domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	devices := make([]*domain.Device, 0)
	for _, device := range r.devices {
		devices = append(devices, device)
	}

	return devices, nil
}

// SignTransaction signs a transaction.
func (r *InMemoryRepository) SignTransaction(deviceID string, data []byte) (*domain.SignatureResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// TODO: Implement signature logic.
	_, exists := r.devices[deviceID]
	if !exists {
		return nil, domain.ErrDeviceNotFound
	}

	return &domain.SignatureResult{
		Signature:  "",
		SignedData: "",
	}, nil
}
