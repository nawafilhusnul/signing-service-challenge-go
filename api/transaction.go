package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func (s *Server) SignTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var req SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid JSON"})
		return
	}

	device, err := s.transactionService.SignTransaction(req.DeviceID, req.Data)
	if err != nil {
		switch err {
		case domain.ErrDeviceNotFound:
			WriteErrorResponse(w, http.StatusNotFound, []string{err.Error()})
		case domain.ErrInvalidDeviceID:
			WriteErrorResponse(w, http.StatusBadRequest, []string{err.Error()})
		default:
			WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}
		return
	}

	WriteAPIResponse(w, http.StatusOK, device)
}
