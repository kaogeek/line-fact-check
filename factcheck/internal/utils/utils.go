// Package utils provides commonly shared code
package utils

import (
	"fmt"

	"github.com/google/uuid"
)

// DefaultIfZero returns v if v is non-zero,
// and falls back to default d otherwise
func DefaultIfZero[T comparable](v, d T) T {
	var zero T
	if v != zero {
		return v
	}
	return d
}

func Ptr[T any](v T) *T {
	return &v
}

func NewID() interface{ String() string } {
	return uuid.New()
}

func Map[T any, U any](slice []T, fn func(T) (U, error)) ([]U, error) {
	if len(slice) == 0 {
		return nil, nil
	}
	result := make([]U, len(slice))
	for i := range slice {
		var err error
		result[i], err = fn(slice[i])
		if err != nil {
			return nil, fmt.Errorf("MapSlice[%d]: %w", i, err)
		}
	}
	return result, nil
}

func MapNoError[T any, U any](slice []T, fn func(T) U) []U {
	if len(slice) == 0 {
		return nil
	}
	result := make([]U, len(slice))
	for i := range slice {
		result[i] = fn(slice[i])
	}
	return result
}

func String[S1 ~string, S2 ~string](s1 S1) S2 {
	return S2(s1)
}
