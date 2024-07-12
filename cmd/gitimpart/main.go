package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mumoshu/gitimpart"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flagset := flag.NewFlagSet("gitimpart", flag.ContinueOnError)

	vars := make(map[string]string)

	file := flagset.String("file", "gitimpart.jsonnet", "The configuration file for rendering and pushing files")
	ghTokenEnv := flagset.String("github-token-env", "GITHUB_TOKEN", "The environment variable name that contains the GitHub token")
	repo := flagset.String("repo", "", "The repository to push the changes to. It should be in the format of `https://github.com/USER/REPO.git`")
	branch := flagset.String("branch", "main", "The branch to push the changes to")
	dryRun := flagset.Bool("dry-run", false, "Print the changes that would be made without actually making them")
	pullRequest := flagset.Bool("pull-request", false, "Send a pull request to the branch after pushing the changes, instead of pushing directly to the branch")

	flagset.Func("var", "The variables to pass to the jsonnet file. Variables are available via std.extVar(name)", func(v string) error {
		fields := strings.Split(v, ",")
		for _, kv := range fields {
			kv := strings.Split(kv, "=")
			if len(kv) != 2 {
				return fmt.Errorf("invalid format for -var: %s", v)
			}
			vars[kv[0]] = kv[1]
		}
		return nil
	})

	if err := flagset.Parse(args); err != nil {
		return fmt.Errorf("failed to parse the flags: %v", err)
	}

	ghtoken := os.Getenv(*ghTokenEnv)
	if ghtoken == "" {
		log.Printf("GITHUB_TOKEN is not set. Access to private repositories will be denied unless you configure other means of authentication")
	}

	var loadOpts []gitimpart.LoadOption

	if len(vars) > 0 {
		loadOpts = append(loadOpts, gitimpart.Vars(vars))
	}

	r, err := gitimpart.RenderFile(*file, loadOpts...)
	if err != nil {
		return fmt.Errorf("failed to render file %s: %v", *file, err)
	}

	opts := []gitimpart.PushOptions{
		gitimpart.WithGitHubToken(ghtoken),
	}

	if *dryRun {
		opts = append(opts, gitimpart.WithDryRun())
	}

	if *pullRequest {
		opts = append(opts, gitimpart.WithPullRequest())
	}

	err = gitimpart.Push(
		*r,
		*repo,
		*branch,
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to push the changes: %v", err)
	}

	if *dryRun {
		log.Printf("successfully pushed the changes to %s/%s", *repo, *branch)
	}

	return nil
}
