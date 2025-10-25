package domain

import "errors"

var (
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDeviceAlreadyExists = errors.New("device already exists")
	ErrInvalidAlgorithm    = errors.New("invalid algorithm")
	ErrInvalidDeviceID     = errors.New("invalid device ID")
	ErrEmptyData           = errors.New("data to sign cannot be empty")
)
