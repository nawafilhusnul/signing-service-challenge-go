package api_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *mux.Router {
	repo := persistence.NewInMemoryRepository()
	svc := service.NewDeviceService(repo)
	srv := api.NewServer("", svc)
	router := mux.NewRouter()
	router.HandleFunc("/api/v0/devices/{deviceId}/sign", srv.SignTransaction).Methods(http.MethodPost)
	router.HandleFunc("/api/v0/devices", srv.CreateDevice).Methods(http.MethodPost)
	router.HandleFunc("/api/v0/devices/{deviceId}", srv.GetDevice).Methods(http.MethodGet)
	router.HandleFunc("/api/v0/devices", srv.GetAllDevices).Methods(http.MethodGet)
	return router
}

func TestServer_CreateDevice(t *testing.T) {
	t.Run("success to register a new device", func(t *testing.T) {
		id := uuid.New().String()
		json := []byte(`{
			"id": "` + id + `",
			"algorithm": "RSA",
			"label": "Device 1"
		}`)
		req, err := http.NewRequest("POST", "/api/v0/devices", bytes.NewReader(json))
		if err != nil {
			t.Fatal(err)
		}

		router := setupTestServer()

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode := rr.Code

		assert.Equal(t, http.StatusCreated, statusCode)

	})

	t.Run("invalid algorithm", func(t *testing.T) {
		id := uuid.New().String()
		json := []byte(`{
			"id": "` + id + `",
			"algorithm": "INVALID_ALG",
			"label": "Device 1"
		}`)
		req, err := http.NewRequest("POST", "/api/v0/devices", bytes.NewReader(json))
		if err != nil {
			t.Fatal(err)
		}

		router := setupTestServer()

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode := rr.Code

		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("empty deviceId", func(t *testing.T) {
		json := []byte(`{
			"algorithm": "RSA",
			"label": "Device 1"
		}`)
		req, err := http.NewRequest("POST", "/api/v0/devices", bytes.NewReader(json))
		assert.NoError(t, err)

		router := setupTestServer()

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		statusCode := rr.Code

		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("duplicate deviceId", func(t *testing.T) {
		id := uuid.New().String()
		json := []byte(`{
			"id": "` + id + `",
			"algorithm": "RSA",
			"label": "Device 1"
		}`)
		req, err := http.NewRequest("POST", "/api/v0/devices", bytes.NewReader(json))
		if err != nil {
			t.Fatal(err)
		}

		router := setupTestServer()

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode := rr.Code

		assert.Equal(t, http.StatusCreated, statusCode)

		req, err = http.NewRequest("POST", "/api/v0/devices", bytes.NewReader(json))
		assert.NoError(t, err)

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode = rr.Code

		assert.Equal(t, http.StatusConflict, statusCode)
	})
}

func TestServer_SignTransaction(t *testing.T) {
	t.Run("success to sign a transaction", func(t *testing.T) {
		router := setupTestServer()

		id := uuid.New().String()
		json := []byte(`{
			"id": "` + id + `",
			"algorithm": "RSA",
			"label": "Device 1"
		}`)
		req, err := http.NewRequest("POST", "/api/v0/devices", bytes.NewReader(json))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode := rr.Code

		assert.Equal(t, http.StatusCreated, statusCode)

		req, err = http.NewRequest("POST", fmt.Sprintf("/api/v0/devices/%s/sign", id), bytes.NewReader([]byte(`{"data": "COFFEE:20251026"}`)))
		assert.NoError(t, err)

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode = rr.Code

		assert.Equal(t, http.StatusOK, statusCode)
	})
	t.Run("failed to sign a transaction, device not found", func(t *testing.T) {
		router := setupTestServer()

		id := uuid.New().String()
		json := []byte(`{
			"id": "` + id + `",
			"algorithm": "RSA",
			"label": "Device 1"
		}`)
		req, err := http.NewRequest("POST", "/api/v0/devices", bytes.NewReader(json))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode := rr.Code

		assert.Equal(t, http.StatusCreated, statusCode)

		req, err = http.NewRequest("POST", fmt.Sprintf("/api/v0/devices/%s/sign", uuid.New().String()), bytes.NewReader([]byte(`{"data": "COFFEE:20251026"}`)))
		assert.NoError(t, err)

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode = rr.Code

		assert.NotEqual(t, http.StatusOK, statusCode)
	})
	t.Run("failed to sign a transaction, missing deviceId", func(t *testing.T) {
		router := setupTestServer()

		req, err := http.NewRequest("POST", "/api/v0/devices/sign", bytes.NewReader([]byte(`{"data": "COFFEE:20251026"}`)))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode := rr.Code

		assert.NotEqual(t, http.StatusOK, statusCode)
	})
}

func TestServer_GetDevice(t *testing.T) {
	t.Run("success to get a device", func(t *testing.T) {
		router := setupTestServer()

		id := uuid.New().String()
		json := []byte(`{
			"id": "` + id + `",
			"algorithm": "RSA",
			"label": "Device 1"
		}`)
		req, err := http.NewRequest("POST", "/api/v0/devices", bytes.NewReader(json))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode := rr.Code

		assert.Equal(t, http.StatusCreated, statusCode)

		req, err = http.NewRequest("GET", fmt.Sprintf("/api/v0/devices/%s", id), nil)
		assert.NoError(t, err)

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode = rr.Code

		assert.Equal(t, http.StatusOK, statusCode)
	})
	t.Run("failed to get a device, device not found", func(t *testing.T) {
		router := setupTestServer()

		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v0/devices/%s", uuid.New().String()), nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode := rr.Code

		assert.Equal(t, http.StatusNotFound, statusCode)
	})
}

func TestServer_GetAllDevices(t *testing.T) {
	t.Run("success to get all devices", func(t *testing.T) {
		router := setupTestServer()

		id := uuid.New().String()
		json := []byte(`{
			"id": "` + id + `",
			"algorithm": "RSA",
			"label": "Device 1"
		}`)
		req, err := http.NewRequest("POST", "/api/v0/devices", bytes.NewReader(json))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode := rr.Code

		assert.Equal(t, http.StatusCreated, statusCode)

		req, err = http.NewRequest("GET", "/api/v0/devices", nil)
		assert.NoError(t, err)

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		statusCode = rr.Code

		assert.Equal(t, http.StatusOK, statusCode)
	})
}
