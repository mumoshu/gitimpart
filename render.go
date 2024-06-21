package gitimpart

import (
	"encoding/json"
)

type Contents struct {
	Files map[string]interface{} `json:"$files"`
}

// RenderFile loads contents from a json or jsonnet file and returns it as a Contents struct.
func RenderFile(path string) (*Contents, error) {
	file, err := LoadFile(path)
	if err != nil {
		return nil, err
	}

	var c Contents

	if err := json.Unmarshal(file, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
