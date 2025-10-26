package service_test

import (
	"encoding/base64"
	"fmt"
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_deviceService_CreateDevice(t *testing.T) {
	t.Run("create device", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := service.NewDeviceService(repository)

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
		deviceService := service.NewDeviceService(repository)

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

func Test_deviceService_SignTransaction(t *testing.T) {
	t.Run("sign single transaction", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := service.NewDeviceService(repository)

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

		// sign a transaction
		trxData := "COFFEE:2025-10-26T07:00:00Z"
		signData := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, trxData, device.LastSignature)

		result, err := deviceService.SignTransaction(id, trxData)

		assert.NoError(t, err, "should not fail to sign transaction")
		assert.Equal(t, signData, result.SignedData)
		assert.Equal(t, 1, device.SignatureCounter)
	})

	t.Run("sign consecutive transactions", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := service.NewDeviceService(repository)

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

		// sign a transaction
		trxData := "COFFEE:2025-10-26T07:00:00Z"
		signData := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, trxData, device.LastSignature)

		result, err := deviceService.SignTransaction(id, trxData)

		assert.NoError(t, err, "should not fail to sign transaction")
		assert.Equal(t, signData, result.SignedData)
		assert.Equal(t, 1, device.SignatureCounter)

		// sign another transaction
		trxData = "COFFEE:2025-10-26T07:01:00Z"
		signData = fmt.Sprintf("%d_%s_%s", device.SignatureCounter, trxData, device.LastSignature)

		result, err = deviceService.SignTransaction(id, trxData)

		assert.NoError(t, err, "should not fail to sign transaction")
		assert.Equal(t, signData, result.SignedData)
		assert.Equal(t, 2, device.SignatureCounter)
	})

	t.Run("sign transaction with invalid device", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := service.NewDeviceService(repository)

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

		// sign a transaction
		trxData := "COFFEE:2025-10-26T07:00:00Z"
		signData := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, trxData, device.LastSignature)

		result, err := deviceService.SignTransaction(id, trxData)

		assert.NoError(t, err, "should not fail to sign transaction")
		assert.Equal(t, signData, result.SignedData)
		assert.Equal(t, 1, device.SignatureCounter)
	})

	t.Run("sign transaction concurrently", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := service.NewDeviceService(repository)

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

		numConcurrent := 10

		wg := sync.WaitGroup{}
		for i := 0; i < numConcurrent; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				// sign a transaction
				trxData := fmt.Sprintf("COFFEE%d:2025-10-26T07:00:00Z", idx)
				_, err := deviceService.SignTransaction(id, trxData)
				assert.NoError(t, err, "should not fail to sign transaction")
			}(i)
		}
		wg.Wait()

		updatedDevice, err := repository.GetByID(id)
		assert.NoError(t, err, "should be able to fetch device")

		assert.Equal(t, numConcurrent, updatedDevice.SignatureCounter,
			fmt.Sprintf("counter should be %d after %d concurrent transactions", numConcurrent, numConcurrent))
	})
}

func Test_deviceService_FindAll(t *testing.T) {
	// spawn repository
	repository := persistence.NewInMemoryRepository()

	// spawn device service
	deviceService := service.NewDeviceService(repository)

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
	// create another device
	id2 := uuid.New().String()
	device2 := &domain.Device{
		ID:               id2,
		Algorithm:        "ECC",
		Label:            "device-2",
		PrivateKey:       "",
		PublicKey:        "",
		SignatureCounter: 0,
		LastSignature:    "",
	}

	err = deviceService.CreateDevice(device2)
	assert.NoError(t, err, "should not fail to create device")

	// find all devices
	devices, err := deviceService.FindAll()
	assert.NoError(t, err, "should not fail to find all devices")
	assert.Equal(t, 2, len(devices), "should find 2 devices")
}

func Test_deviceService_GetDevice(t *testing.T) {
	// spawn repository
	repository := persistence.NewInMemoryRepository()

	// spawn device service
	deviceService := service.NewDeviceService(repository)

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

	// get device
	gotDevice, err := deviceService.GetDevice(id)
	assert.NoError(t, err, "should not fail to get device")
	assert.Equal(t, id, gotDevice.ID)

	decodedID, err := base64.RawStdEncoding.DecodeString(gotDevice.LastSignature)
	assert.NoError(t, err, "should not fail to decode device last signature")
	assert.Equal(t, id, string(decodedID), "device last signature should match device id")
}
