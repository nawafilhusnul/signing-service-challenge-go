package service_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_transactionService_SignTransaction(t *testing.T) {
	t.Run("sign single transaction", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := service.NewDeviceService(repository)

		// spawn transaction service
		transactionService := service.NewTransactionService(repository)

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

		result, err := transactionService.SignTransaction(id, trxData)

		assert.NoError(t, err, "should not fail to sign transaction")
		assert.Equal(t, signData, result.SignedData)
		assert.Equal(t, 1, device.SignatureCounter)
	})

	t.Run("sign consecutive transactions", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := service.NewDeviceService(repository)

		// spawn transaction service
		transactionService := service.NewTransactionService(repository)

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

		result, err := transactionService.SignTransaction(id, trxData)

		assert.NoError(t, err, "should not fail to sign transaction")
		assert.Equal(t, signData, result.SignedData)
		assert.Equal(t, 1, device.SignatureCounter)

		// sign another transaction
		trxData = "COFFEE:2025-10-26T07:01:00Z"
		signData = fmt.Sprintf("%d_%s_%s", device.SignatureCounter, trxData, device.LastSignature)

		result, err = transactionService.SignTransaction(id, trxData)

		assert.NoError(t, err, "should not fail to sign transaction")
		assert.Equal(t, signData, result.SignedData)
		assert.Equal(t, 2, device.SignatureCounter)
	})

	t.Run("sign transaction with invalid device", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := service.NewDeviceService(repository)

		// spawn transaction service
		transactionService := service.NewTransactionService(repository)

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

		result, err := transactionService.SignTransaction(id, trxData)

		assert.NoError(t, err, "should not fail to sign transaction")
		assert.Equal(t, signData, result.SignedData)
		assert.Equal(t, 1, device.SignatureCounter)
	})

	t.Run("sign transaction concurrently", func(t *testing.T) {
		// spawn repository
		repository := persistence.NewInMemoryRepository()

		// spawn device service
		deviceService := service.NewDeviceService(repository)

		// spawn transaction service
		transactionService := service.NewTransactionService(repository)

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
				_, err := transactionService.SignTransaction(id, trxData)
				assert.NoError(t, err, "should not fail to sign transaction")
			}(i)
		}
		wg.Wait()

		assert.Equal(t, numConcurrent, device.SignatureCounter, "should not fail to sign transaction")
	})
}
