// Store is an interface for storing configs.
// There are two implementations of this interface:
// - Local
// - Git
// - PullRequest
package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/mumoshu/gitimpart/config"
	"github.com/mumoshu/gitimpart/convention"
	"github.com/mumoshu/gitimpart/envvar"
)

const appName = "gitimpart"

type Store interface {
	Put(context context.Context, path string, content string) error
	List(context context.Context, path string) ([]string, error)
	Get(context context.Context, path string) (*string, error)
	Delete(context context.Context, path string) error

	// Transact runs the given function within the directory that the store stores the configs.
	//
	// The path argument is the path to the directory.
	// It can be a temporary directory or the specified directory in the clone of the gitops repository,
	// or whatever the store implementation relies on.
	//
	// The function is expected to return the result of the rendering, which contains
	// added, modified, and deleted files.
	// The function is expected to return an error if the rendering fails.
	//
	// The store implementation is expected to include the files returned by the function
	// in the commit.
	// The caller is expected to call Commit after calling Transact.
	Transact(fn func(path string) (*RenderResult, error)) (*RenderResult, error)

	// Commit commits the changes made to the store.
	// The subject and body are used as the commit message, if applicable.
	// If the store does not support commits, it returns nil.
	Commit(context context.Context, subject, body string) error
}

// Make makes a store based on the config.Delegate.
func Make(id string, t time.Time, d *config.Delegate) Store {
	if d == nil {
		return newLocal(id)
	}

	repoURL := convention.RepoURL(d.Git.Repo)
	g := newGit(id, t, d)

	if d.PullRequest != nil {
		return &PullRequest{
			RepositoryURL: repoURL,
			Git:           g,
		}
	}

	return g
}

func newGit(id string, t time.Time, d *config.Delegate) *Git {
	baseBranch := os.Getenv(envvar.BaseBranch)
	if d.Git.Branch != "" {
		baseBranch = d.Git.Branch
	}

	auth := &http.BasicAuth{
		Username: appName + "bot", // This can be anything except an empty string
		Password: os.Getenv(envvar.GitHubToken),
	}

	var newBranch string

	if d.PullRequest != nil {
		newBranch = fmt.Sprintf(appName+"/%s-%s", id, t.Format("20060102150405"))
	}

	gitRoot := os.Getenv(envvar.GitRoot)
	if gitRoot == "" {
		gitRoot = "." + appName + "/repositories"
	}

	repoURL := convention.RepoURL(d.Git.Repo)

	g := NewGit(
		auth,
		baseBranch,
		newBranch,
		repoURL,
		os.Getenv(envvar.GitCommitAuthorUserName),
		os.Getenv(envvar.GitCommitAuthorEmail),
		gitRoot,
		d.Git.Push,
	)

	return g
}
