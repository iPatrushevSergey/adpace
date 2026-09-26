package port

import "github.com/google/uuid"

type IDGenerator interface {
	NewID() (uuid.UUID, error)
}
