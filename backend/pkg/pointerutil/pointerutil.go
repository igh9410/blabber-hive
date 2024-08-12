package pointerutil

import (
	"time"

	"github.com/google/uuid"
)

// StringToPtr converts a string to a *string.
func StringToPtr(s string) *string {
	return &s
}

// PtrToString converts a *string to a string. If the pointer is nil, it returns an empty string.
func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Int32ToPtr converts an int32 to a *int32.
func Int32ToPtr(i int32) *int32 {
	return &i
}

// PtrToInt32 converts a *int32 to an int32. If the pointer is nil, it returns 0.
func PtrToInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

// Int64ToPtr converts an int64 to a *int64.
func Int64ToPtr(i int64) *int64 {
	return &i
}

// PtrToInt64 converts a *int64 to an int64. If the pointer is nil, it returns 0.
func PtrToInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// BoolToPtr converts a bool to a *bool.
func BoolToPtr(b bool) *bool {
	return &b
}

// PtrToBool converts a *bool to a bool. If the pointer is nil, it returns false.
func PtrToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Float64ToPtr converts a float64 to a *float64.
func Float64ToPtr(f float64) *float64 {
	return &f
}

// PtrToFloat64 converts a *float64 to a float64. If the pointer is nil, it returns 0.0.
func PtrToFloat64(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

// UUIDToPtr converts a uuid.UUID to a *uuid.UUID.
func UUIDToPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

// PtrToUUID converts a *uuid.UUID to a uuid.UUID. If the pointer is nil, it returns a zero uuid.UUID.
func PtrToUUID(u *uuid.UUID) uuid.UUID {
	if u == nil {
		return uuid.UUID{}
	}
	return *u
}

// TimeToPtr converts a time.Time to a *time.Time.
func TimeToPtr(t time.Time) *time.Time {
	return &t
}

// PtrToTime converts a *time.Time to a time.Time. If the pointer is nil, it returns a zero time.Time.
func PtrToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
