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

	err := s.repository.Create(&domain.Device{
		ID:               req.ID,
		Algorithm:        req.Algorithm,
		Label:            req.Label,
		SignatureCounter: 0,
		LastSignature:    "",
		PrivateKey:       nil,
		PublicKey:        nil,
		CreatedAt:        time.Now(),
	})
	if err != nil {
		switch err {
		case domain.ErrInvalidAlgorithm, domain.ErrInvalidDeviceID, domain.ErrDeviceAlreadyExists:
			WriteErrorResponse(w, http.StatusBadRequest, []string{err.Error()})
		default:
			WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}
		return
	}

	WriteAPIResponse(w, http.StatusCreated, nil)
}
