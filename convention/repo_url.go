package convention

import (
	"fmt"
	"os"
	"strings"

	"github.com/mumoshu/gitimpart/envvar"
)

// RepoURL returns the normalized URL of the repository.
func RepoURL(repo string) string {
	var repoURL string
	if strings.Count(repo, "/") == 1 {
		githubBaseURL := "https://github.com/"
		if os.Getenv(envvar.GitHubEnterpriseURL) != "" {
			githubBaseURL = os.Getenv(envvar.GitHubEnterpriseURL)
		}
		repoURL = githubBaseURL + repo + ".git"
	} else if strings.Count(repo, "/") == 2 {
		repoURL = "https://" + repo + ".git"
	} else if strings.HasPrefix(repo, "https://") {
		repoURL = repo
	} else {
		panic(fmt.Sprintf("invalid repo: %s", repo))
	}
	return repoURL
}
