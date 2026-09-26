package generator

import "github.com/google/uuid"

type IDGenerator struct{}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (g *IDGenerator) NewID() (uuid.UUID, error) {
	return uuid.NewV7()
}
