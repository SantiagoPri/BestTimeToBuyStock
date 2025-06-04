package stock

import "backend/pkg/errors"

var (
	ErrNotFound = errors.New(errors.ErrNotFound, "stock not found")
)
