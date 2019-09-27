package db

import "github.com/twinj/uuid"

//UUIDer generates UUIDs
//go:generate counterfeiter . UUIDer
type UUIDer interface {
	UUID() string
}

//UUIDGenerator is a standard implementation of UUIDer
type UUIDGenerator struct{}

//UUID generates a new V4 UUID
func (UUIDGenerator) UUID() string {
	return uuid.NewV4().String()
}
