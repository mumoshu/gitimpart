package gitimpart

import (
	"encoding/json"
	"path/filepath"
)

type Contents struct {
	Files     map[string]interface{}            `json:"$files"`
	Kustomize map[string]map[string]interface{} `json:"$kustomize"`
}

type LoadConfig struct {
	Vars map[string]string
}

type LoadOption func(*LoadConfig)

func Vars(vars map[string]string) LoadOption {
	return func(c *LoadConfig) {
		c.Vars = vars
	}
}

// RenderFile loads contents from a json or jsonnet file and returns it as a Contents struct.
func RenderFile(path string, opts ...LoadOption) (*Contents, error) {
	file, err := LoadFile(path, opts...)
	if err != nil {
		return nil, err
	}

	var c Contents

	if err := json.Unmarshal(file, &c); err != nil {
		return nil, err
	}

	for dir, files := range c.Kustomize {
		for name, content := range files {
			c.Files[filepath.Join(dir, name)] = content
			// We no longer need the file content as
			// its content is already in c.Files.
			files[name] = nil
		}
	}

	return &c, nil
}
