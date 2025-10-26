package service

import (
	"encoding/base64"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

type DeviceService interface {
	CreateDevice(device *domain.Device) error
	GetDevice(deviceID string) (*domain.Device, error)
	FindAll() ([]*domain.Device, error)
	SignTransaction(deviceID string, data string) (*domain.SignatureResult, error)
}

type deviceService struct {
	repository persistence.Repository
}

func NewDeviceService(repository persistence.Repository) DeviceService {
	return &deviceService{repository: repository}
}

func (s *deviceService) CreateDevice(device *domain.Device) error {
	lastSignature := base64.RawStdEncoding.EncodeToString([]byte(device.ID))
	gen, err := crypto.NewGenerator(device.Algorithm)
	if err != nil {
		return err
	}

	keyPair, err := gen.Generate()
	if err != nil {
		return err
	}

	device.PrivateKey = string(keyPair.GetPrivateKeyPEM())
	device.PublicKey = string(keyPair.GetPublicKeyPEM())

	// pre-sign the device
	device.SignatureCounter = 0
	device.LastSignature = lastSignature

	return s.repository.Create(device)
}

func (s *deviceService) SignTransaction(deviceID string, data string) (*domain.SignatureResult, error) {
	var result *domain.SignatureResult

	_, err := s.repository.Update(deviceID, func(device *domain.Device) error {
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

func (s *deviceService) GetDevice(deviceID string) (*domain.Device, error) {
	return s.repository.GetByID(deviceID)
}

func (s *deviceService) FindAll() ([]*domain.Device, error) {
	return s.repository.FindAll()
}
