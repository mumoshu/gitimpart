package gitimpart

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
)

// LoadFile loads a json or jsonnet file and returns the content as a byte slice.
// In case it is a jsonnet file, it evaluates the jsonnet file and returns the resulting json as a byte slice.
func LoadFile(path string, opts ...LoadOption) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if filepath.Ext(path) == ".jsonnet" {
		vm := jsonnet.MakeVM()

		if strings.Contains(filepath.Base(path), ".template.") {
			vm.ExtVar("template", "true")

			if v := os.Getenv("GITHUB_REPOSITORY"); v != "" {
				splits := strings.Split(v, "/")
				vm.ExtVar("github_repo_owner", splits[0])
				vm.ExtVar("github_repo_name", splits[1])
			} else {
				return nil, errors.New(".template.jsonnet requires GITHUB_REPOSITORY to be set to OWNER/REPO_NAME for the template to access `std.extVar(\"github_repo_name\")` and `std.extVar(\"github_repo_owner\")`")
			}
		}

		json, err := vm.EvaluateAnonymousSnippet(path, string(file))
		if err != nil {
			return nil, err
		}

		file = []byte(json)
	}

	return file, nil
}
