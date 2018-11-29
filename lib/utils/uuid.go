package utils

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

func NewUUID() string {
	id, err := uuid.NewV4()
	if err != nil {
		return fmt.Sprintf("%v", time.Now().Nanosecond())
	}
	return id.String()
}
