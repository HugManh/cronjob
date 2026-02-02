package dto_tasks

import (
	"fmt"
	"strings"
)

// BoolString is a custom type that can unmarshal JSON boolean values
type BoolString bool

func (b *BoolString) UnmarshalJSON(data []byte) error {

	str := strings.Trim(string(data), "\"")
	switch str {
	case "true", "True", "1":
		*b = true
	case "false", "False", "0":
		*b = false
	case "":
		return nil
	default:
		return fmt.Errorf("invalid boolean: %s", str)
	}
	return nil
}

type AddTaskRequest struct {
	Name    string      `json:"name"`
	Execute string      `json:"execute"` // vd: "*/10 * * * * *"
	Message string      `json:"message"`
	Active  *BoolString `json:"active"`
}
