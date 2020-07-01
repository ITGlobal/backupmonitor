package model

import (
	"encoding/json"
)

// Empty is an empty JSON object
type Empty struct {
}

func toJSON(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}
