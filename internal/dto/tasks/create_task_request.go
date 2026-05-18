package tasks

import (
	"fmt"
	"strings"
)

type BoolString bool

func (b *BoolString) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), "\"")
	switch value {
	case "true", "True", "1":
		*b = true
	case "false", "False", "0":
		*b = false
	case "":
		return nil
	default:
		return fmt.Errorf("invalid boolean: %s", value)
	}
	return nil
}

type CreateTaskRequest struct {
	Name    string      `json:"name"`
	Execute string      `json:"execute"`
	Message string      `json:"message"`
	Active  *BoolString `json:"active"`
}
