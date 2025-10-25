package service

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

type DeviceService interface {
	CreateDevice(device *domain.Device) error
}

type deviceService struct {
	repository persistence.Repository
}

func NewDeviceService(repository persistence.Repository) DeviceService {
	return &deviceService{repository: repository}
}

func (s *deviceService) CreateDevice(device *domain.Device) error {
	gen, err := crypto.NewGenerator(device.Algorithm)
	if err != nil {
		return err
	}

	keyPair, err := gen.Generate()
	if err != nil {
		return err
	}

	device.PrivateKey = keyPair.GetPrivateKeyPEM()
	device.PublicKey = keyPair.GetPublicKeyPEM()

	return s.repository.Create(device)
}
