// Package utils provides commonly shared code
package utils

// DefaultIfZero returns v if v is non-zero,
// and falls back to default d otherwise
func DefaultIfZero[T comparable](v, d T) T {
	var zero T
	if v != zero {
		return v
	}
	return d
}
