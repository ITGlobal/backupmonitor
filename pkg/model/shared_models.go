package model

import (
	"encoding/json"
)

func toJSON(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}
