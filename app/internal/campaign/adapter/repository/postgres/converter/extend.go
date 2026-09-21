package converter

import (
	"time"

	"github.com/google/uuid"
)

// UUIDToString converts a database UUID to a domain id.
func UUIDToString(id uuid.UUID) string {
	return id.String()
}

// StringToUUID converts a domain id to a database UUID.
func StringToUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// CopyTime maps time.Time for goverter (required: time.Time has unexported fields).
func CopyTime(v time.Time) time.Time {
	return v
}

// CopyTimePtr maps *time.Time for goverter (required: time.Time has unexported fields).
func CopyTimePtr(v *time.Time) *time.Time {
	return v
}
