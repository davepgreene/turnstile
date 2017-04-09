package utils

import (
	"github.com/satori/go.uuid"
)

// CreateCorrelationID generates a new uuid
func CreateCorrelationID() string {
	return uuid.NewV4().String()
}
