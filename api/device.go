package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func (s *Server) CreateDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var req CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid JSON"})
		return
	}

	newDevice := domain.Device{
		ID:        req.ID,
		Algorithm: req.Algorithm,
		Label:     req.Label,
		CreatedAt: time.Now(),
	}

	err := s.deviceService.CreateDevice(&newDevice)
	if err != nil {
		switch err {
		case domain.ErrDeviceAlreadyExists:
			WriteErrorResponse(w, http.StatusConflict, []string{err.Error()})
		case domain.ErrInvalidAlgorithm, domain.ErrInvalidDeviceID:
			WriteErrorResponse(w, http.StatusBadRequest, []string{err.Error()})
		default:
			WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}
		return
	}

	WriteAPIResponse(w, http.StatusCreated, newDevice)
}
