package utils

import "encoding/json"

func IsJson(s string) error {
	var js struct{}

	// o json.Unmarshal tenta converter a string para um json
	if err := json.Unmarshal([]byte(s), &js); err != nil {
		return err
	}

	return nil
}
