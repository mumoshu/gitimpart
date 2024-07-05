package gitimpart

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/mumoshu/gitimpart/store"
	"gopkg.in/yaml.v2"
)

type PushConfig struct {
	// Auth is the authentication information to use when committing the changes
	Auth *http.BasicAuth
	// Dir is the directory to store the git repository.
	// If provided, the caller needs to clean up the directory after the commit.
	Dir string
	// Subject is the commit message subject.
	Subject string
	// Body is the commit message body.
	Body string
	// SendPullRequest is a flag to send a pull request after the commit-push.
	SendPullRequest bool
	// KustomizeBin is the path to the kustomize binary.
	KustomizeBin string
}

type PushOptions func(*PushConfig)

func WithGitHubToken(token string) PushOptions {
	return func(c *PushConfig) {
		c.Auth = &http.BasicAuth{
			Username: "gitimpartbot",
			Password: token,
		}
	}
}

func WithPullRequest() PushOptions {
	return func(c *PushConfig) {
		c.SendPullRequest = true
	}
}

func WithKustomizeBin(bin string) PushOptions {
	return func(c *PushConfig) {
		c.KustomizeBin = bin
	}
}

// Push pushes the contents to the specified repository and branch.
//
// When the Dir field is not provided, it creates a temporary directory to store the git repository.
// When the Dir field is provided, the caller needs to clean up the directory after the commit.
//
// The Auth field is required. Set a valid GitHub token to the Password field.
//
// The Subject field is the commit message subject. If not provided, it uses the branch name.
func Push(r Contents, repo, branch string, opts ...PushOptions) error {
	var c PushConfig
	for _, o := range opts {
		o(&c)
	}

	if c.Auth.Password == "" {
		return fmt.Errorf("CommitConfig.Auth.Password is required. Set a valid GitHub token to CommitConfig.Auth.Password")
	}

	tm := time.Now()

	name := tm.Format("20060102150405")
	newBranch := fmt.Sprintf("gitimpart-%s", name)

	dir := c.Dir
	if dir == "" {
		var err error
		dir, err = os.MkdirTemp("", newBranch)
		if err != nil {
			return fmt.Errorf("unable to create temp dir: %w", err)
		}

		defer os.RemoveAll(dir)
	}

	gitRoot := filepath.Join(dir, ".gitimpart", "gitroot")
	if err := os.MkdirAll(gitRoot, 0755); err != nil {
		return fmt.Errorf("unable to create git root: %w", err)
	}

	var s store.Store

	g := store.NewGit(
		c.Auth,
		branch,
		newBranch,
		repo,
		"test author", "test@example.com",
		gitRoot,
		true,
	)

	if c.SendPullRequest {
		s = &store.PullRequest{
			RepositoryURL: repo,
			Git:           g,
		}
	} else {
		s = g
	}

	_, err := s.Transact(func(dir string) (*store.RenderResult, error) {
		var updates []string
		for name, content := range r.Files {
			p := filepath.Join(dir, name)
			if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
				return nil, fmt.Errorf("mkdir error: %w", err)
			}

			switch content := content.(type) {
			case string:
				if err := os.WriteFile(p, []byte(content), 0644); err != nil {
					return nil, fmt.Errorf("write error: %w", err)
				}
			default:
				switch filepath.Ext(name) {
				case ".json":
					b, err := json.Marshal(content)
					if err != nil {
						return nil, fmt.Errorf("marshal error: %w", err)
					}
					if err := os.WriteFile(p, b, 0644); err != nil {
						return nil, fmt.Errorf("write error: %w", err)
					}
				case ".yaml", ".yml":
					b, err := yaml.Marshal(content)
					if err != nil {
						return nil, fmt.Errorf("marshal error: %w", err)
					}
					if err := os.WriteFile(p, b, 0644); err != nil {
						return nil, fmt.Errorf("write error: %w", err)
					}
				default:
					return nil, fmt.Errorf("unsupported file type: %s", name)
				}
			}

			updates = append(updates, name)
		}

		var kustomizeBinAbs string

		for kDir, files := range r.Kustomize {
			if kustomizeBinAbs == "" {
				kustomizeBin := "kustomize"
				if c.KustomizeBin != "" {
					kustomizeBin = c.KustomizeBin
				}
				kustomizeBin, err := filepath.Abs(kustomizeBin)
				if err != nil {
					return nil, fmt.Errorf("unable to get absolute path to kustomize binary %s: %w", kustomizeBin, err)
				}
				kustomizeBinAbs, err = exec.LookPath(kustomizeBin)
				if err != nil {
					return nil, fmt.Errorf("unable to find kustomize binary: %w", err)
				}
			}

			thisDir := filepath.Join(dir, kDir)

			// Executes kustomize edit add resource for each file in files and adds the kustomization.yaml as updated
			// file.
			for name := range files {
				// Create the directory that should contain the kustomization.yaml file,
				// if it does not exist.
				if stat, err := os.Stat(thisDir); err != nil {
					if os.IsNotExist(err) {
						if err := os.MkdirAll(thisDir, 0755); err != nil {
							return nil, fmt.Errorf("mkdir error: %w", err)
						}
					} else {
						return nil, fmt.Errorf("stat error: %w", err)
					}
				} else if !stat.IsDir() {
					return nil, fmt.Errorf("not a directory: %s", thisDir)
				}

				// Create the kustomization.yaml file, if it does not exist.
				// Otherwise kustomize-edit-add-resource fails with:
				//   Error: Missing kustomization file 'kustomization.yaml'.
				if _, err := os.Stat(filepath.Join(thisDir, "kustomization.yaml")); err != nil {
					if os.IsNotExist(err) {
						if err := os.WriteFile(filepath.Join(thisDir, "kustomization.yaml"), []byte("resources:\n"), 0644); err != nil {
							return nil, fmt.Errorf("write error: %w", err)
						}
					} else {
						return nil, fmt.Errorf("stat error: %w", err)
					}
				}

				kustomizeEdit := exec.Command(kustomizeBinAbs, "edit", "add", "resource", name)
				kustomizeEdit.Dir = thisDir
				combined, err := kustomizeEdit.CombinedOutput()
				if err != nil {
					return nil, fmt.Errorf("kustomize edit error: %w, %s", err, combined)
				}
			}
			updates = append(updates, filepath.Join(kDir, "kustomization.yaml"))
		}

		return &store.RenderResult{
			AddedOrModifiedFiles: updates,
		}, nil
	})
	if err != nil {
		return fmt.Errorf("unable to transact: %w", err)
	}

	ctx := context.Background()
	subject := c.Subject
	if subject == "" {
		subject = newBranch
	}
	body := c.Body
	if body == "" {
		body = "test"
	}
	if err := s.Commit(ctx, subject, body); err != nil {
		return fmt.Errorf("unable to commit: %w", err)
	}

	return nil
}
