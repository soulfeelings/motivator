package model

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrForbidden       = errors.New("forbidden")
	ErrConflict        = errors.New("already exists")
	ErrDuplicateInvite = errors.New("invite already sent to this email")
	ErrInviteExpired   = errors.New("invite has expired")
	ErrEmailMismatch   = errors.New("email does not match invite")
)
