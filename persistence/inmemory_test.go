package persistence_test

import (
	"testing"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository_Create(t *testing.T) {
	t.Run("sequential creates", func(t *testing.T) {
		tests := []struct {
			name          string
			device        *domain.Device
			wantErr       bool
			expectedError error
		}{
			{
				name: "Create device",
				device: &domain.Device{
					ID:               "1",
					Algorithm:        "RSA",
					Label:            "Test Device",
					SignatureCounter: 0,
					LastSignature:    "",
					PrivateKey:       nil,
					PublicKey:        nil,
					CreatedAt:        time.Now(),
				},
				wantErr:       false,
				expectedError: nil,
			},
			{
				name: "Create duplicate device",
				device: &domain.Device{
					ID:               "1",
					Algorithm:        "ECC",
					Label:            "Duplicate Device",
					SignatureCounter: 0,
					LastSignature:    "",
					PrivateKey:       nil,
					PublicKey:        nil,
					CreatedAt:        time.Now(),
				},
				wantErr:       true,
				expectedError: domain.ErrDeviceAlreadyExists,
			},
		}

		r := persistence.NewInMemoryRepository()
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				gotErr := r.Create(tt.device)
				if tt.wantErr {
					assert.Error(t, gotErr)
					assert.EqualError(t, gotErr, tt.expectedError.Error())
					return
				}

				assert.NoError(t, gotErr)
			})
		}
	})

	t.Run("concurrent creates", func(t *testing.T) {
		r := persistence.NewInMemoryRepository()

		t.Run("different devices", func(t *testing.T) {
			const numConcurrent = 10
			done := make(chan error, numConcurrent)
			successCount := 0

			for i := 0; i < numConcurrent; i++ {
				go func(id int) {
					device := &domain.Device{
						ID:               string(rune('A' + id)),
						Algorithm:        "RSA",
						Label:            "Concurrent Device",
						SignatureCounter: 0,
						LastSignature:    "",
						PrivateKey:       nil,
						PublicKey:        nil,
						CreatedAt:        time.Now(),
					}
					done <- r.Create(device)
				}(i)
			}

			for i := 0; i < numConcurrent; i++ {
				err := <-done
				if err != nil {
					t.Errorf("Concurrent Create() failed: %v", err)
				} else {
					successCount++
				}
			}

			assert.Equal(t, numConcurrent, successCount)
		})

		t.Run("same device", func(t *testing.T) {
			const numConcurrent = 10
			done := make(chan error, numConcurrent)
			successCount := 0

			for i := 0; i < numConcurrent; i++ {
				go func() {
					device := &domain.Device{
						ID:               "concurrent-test",
						Algorithm:        "RSA",
						Label:            "Same Device",
						SignatureCounter: 0,
						LastSignature:    "",
						PrivateKey:       nil,
						PublicKey:        nil,
						CreatedAt:        time.Now(),
					}
					done <- r.Create(device)
				}()
			}

			for i := 0; i < numConcurrent; i++ {
				err := <-done
				if err == nil {
					successCount++
				}
			}

			assert.Equal(t, 1, successCount)
		})
	})
}

func TestInMemoryRepository_GetByID(t *testing.T) {
	t.Run("get existing device", func(t *testing.T) {
		r := persistence.NewInMemoryRepository()

		// prepare a device to be asserted
		device := &domain.Device{
			ID:               "1",
			Algorithm:        "RSA",
			Label:            "Test Device",
			SignatureCounter: 0,
			LastSignature:    "",
			PrivateKey:       nil,
			PublicKey:        nil,
			CreatedAt:        time.Now(),
		}
		err := r.Create(device)
		assert.NoError(t, err)

		got, gotErr := r.GetByID(device.ID)
		assert.NoError(t, gotErr)

		assert.Equal(t, got, device)
	})

	t.Run("get non existing device", func(t *testing.T) {
		r := persistence.NewInMemoryRepository()
		_, gotErr := r.GetByID("1")
		assert.Error(t, gotErr)
		assert.EqualError(t, gotErr, domain.ErrDeviceNotFound.Error())

	})
}

func TestInMemoryRepository_FindAll(t *testing.T) {
	t.Run("find all devices", func(t *testing.T) {
		r := persistence.NewInMemoryRepository()

		// prepare a device to be asserted
		devices := []*domain.Device{
			{
				ID:               "1",
				Algorithm:        "RSA",
				Label:            "Test Device",
				SignatureCounter: 0,
				LastSignature:    "",
				PrivateKey:       nil,
				PublicKey:        nil,
				CreatedAt:        time.Now(),
			},
			{
				ID:               "2",
				Algorithm:        "RSA",
				Label:            "Test Device",
				SignatureCounter: 0,
				LastSignature:    "",
				PrivateKey:       nil,
				PublicKey:        nil,
				CreatedAt:        time.Now(),
			},
		}

		for _, device := range devices {
			err := r.Create(device)
			assert.NoError(t, err)
		}

		got, gotErr := r.FindAll()
		assert.NoError(t, gotErr)
		assert.Equal(t, len(got), len(devices))
	})
}
