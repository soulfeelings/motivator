package model

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrForbidden        = errors.New("forbidden")
	ErrInsufficientFunds = errors.New("insufficient coins")
	ErrCannotAttackSelf = errors.New("cannot attack your own base")
)
