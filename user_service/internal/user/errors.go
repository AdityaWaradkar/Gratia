package user

import "errors"

var (
	// Auth / authorization
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// Validation
	ErrInvalidUserID = errors.New("invalid user id")
	ErrInvalidInput  = errors.New("invalid input")

	// Donor profile
	ErrDonorProfileExists   = errors.New("donor profile already exists")
	ErrDonorProfileNotFound = errors.New("donor profile not found")

	// NGO profile
	ErrNGOProfileExists      = errors.New("ngo profile already exists")
	ErrNGOProfileNotFound    = errors.New("ngo profile not found")
	ErrNGOAlreadyVerified    = errors.New("ngo already verified")
	ErrInvalidNGODetails     = errors.New("invalid ngo details")
	ErrTargetUserIDRequired  = errors.New("target user id required")
)