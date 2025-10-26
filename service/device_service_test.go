package service

import (
	"encoding/base64"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_deviceService_CreateDevice(t *testing.T) {
	t.Run("create device", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := NewDeviceService(repository)

		// create device
		id := uuid.New().String()
		device := &domain.Device{
			ID:               id,
			Algorithm:        "ECC",
			Label:            "device-1",
			PrivateKey:       "",
			PublicKey:        "",
			SignatureCounter: 0,
			LastSignature:    "",
		}

		err := deviceService.CreateDevice(device)
		assert.NoError(t, err, "should not fail to create device")

		// decode device last signature
		lastSignature, err := base64.RawStdEncoding.DecodeString(device.LastSignature)
		assert.NoError(t, err, "should not fail to decode device last signature")
		assert.Equal(t, id, string(lastSignature), "device last signature should match device id")
	})

	t.Run("create device with invalid algorithm", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := NewDeviceService(repository)

		// create device
		id := uuid.New().String()
		device := &domain.Device{
			ID:               id,
			Algorithm:        "INVALID",
			Label:            "device-1",
			PrivateKey:       "",
			PublicKey:        "",
			SignatureCounter: 0,
			LastSignature:    "",
		}

		err := deviceService.CreateDevice(device)
		assert.Error(t, err, "should fail to create device with invalid algorithm")
	})
}
