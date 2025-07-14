// Package utils provides commonly shared code
package utils

import "github.com/google/uuid"

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

func MapSlice[T any, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}
