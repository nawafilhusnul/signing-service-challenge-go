package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/helper"
	"github.com/gorilla/mux"
)

func (s *Server) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid JSON"})
		return
	}

	errs := make([]string, 0)
	if !helper.IsValidUUID(req.ID) {
		errs = append(errs, "Invalid Device ID. UUID format expected")
	}

	if len(errs) > 0 {
		WriteErrorResponse(w, http.StatusBadRequest, errs)
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

func (s *Server) SignTransaction(w http.ResponseWriter, r *http.Request) {
	deviceId := mux.Vars(r)["deviceId"]
	if !helper.IsValidUUID(deviceId) {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid Device ID. UUID format expected"})
		return
	}

	var req SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid JSON"})
		return
	}

	errs := make([]string, 0)
	if req.Data == "" {
		errs = append(errs, domain.ErrEmptyData.Error())
	}

	if len(errs) > 0 {
		WriteErrorResponse(w, http.StatusBadRequest, errs)
		return
	}

	result, err := s.deviceService.SignTransaction(deviceId, req.Data)
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

	WriteAPIResponse(w, http.StatusOK, result)
}

func (s *Server) GetDevice(w http.ResponseWriter, r *http.Request) {
	deviceId := mux.Vars(r)["deviceId"]
	if !helper.IsValidUUID(deviceId) {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid Device ID. UUID format expected"})
		return
	}

	device, err := s.deviceService.GetDevice(deviceId)
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

func (s *Server) GetAllDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := s.deviceService.FindAll()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	WriteAPIResponse(w, http.StatusOK, devices)
}
