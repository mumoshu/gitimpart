package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	file := flagset.String("file", "gitimpart.jsonnet", "The configuration file for rendering and pushing files")
	ghTokenEnv := flagset.String("github-token-env", "GITHUB_TOKEN", "The environment variable name that contains the GitHub token")
	repo := flagset.String("repo", "", "The repository to push the changes to. It should be in the format of `https://github.com/USER/REPO.git`")
	branch := flagset.String("branch", "main", "The branch to push the changes to")
	dryRun := flagset.Bool("dry-run", false, "Print the changes that would be made without actually making them")
	pullRequest := flagset.Bool("pull-request", false, "Send a pull request to the branch after pushing the changes, instead of pushing directly to the branch")

	if err := flagset.Parse(args); err != nil {
		return fmt.Errorf("failed to parse the flags: %v", err)
	}

	ghtoken := os.Getenv(*ghTokenEnv)
	if ghtoken == "" {
		log.Printf("GITHUB_TOKEN is not set. Access to private repositories will be denied unless you configure other means of authentication")
	}

	r, err := gitimpart.RenderFile(*file)
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
