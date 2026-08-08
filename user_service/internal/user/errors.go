package user

import "errors"

// Standard domain-level sentinel errors for uniform error handling across the user service
var (
    // Authorization boundaries
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")

    // Input and identifier validation
    ErrInvalidUserID = errors.New("invalid user id")
    ErrInvalidInput  = errors.New("invalid input")

    // Donor profile states
    ErrDonorProfileExists   = errors.New("donor profile already exists")
    ErrDonorProfileNotFound = errors.New("donor profile not found")

    // NGO profile and verification states
    ErrNGOProfileExists     = errors.New("ngo profile already exists")
    ErrNGOProfileNotFound   = errors.New("ngo profile not found")
    ErrNGOAlreadyVerified   = errors.New("ngo already verified")
    ErrInvalidNGODetails    = errors.New("invalid ngo details")
    ErrTargetUserIDRequired = errors.New("target user id required")
)