package config

type Delegate struct {
	// Git specifies whether the gitops config is loaded from a git repository.
	Git *Git `yaml:"git,omitempty"`

	// PullRequest specifies whether the gitops config is updated via pull request.
	// If false, gitimpart pushes directly to the branch that contains the gitops config.
	// If true, gitimpart creates a feature branch, pushes to the feature branch, and creates a pull request.
	// To be clear, the Branch field serves as the base branch of the pull request.
	PullRequest *PullRequest `yaml:"pullRequest,omitempty"`
}

type Git struct {
	// Repo is either REPO/NAME or URL of the git repository that contains the gitops config.
	// A gitops config can be either a directory or a file, that contains Kubernetes manifests,
	// kustomize config, or Terraform workspaces.
	//
	// This can point to the same repository that the gitimpart.yaml is in and the pull request is made against,
	// or a different target repository that the repository_dispatch is sent to.
	//
	// Regardless, the gitops config is updated in the repository specified by Repo.
	Repo string `yaml:"repo"`

	// Branch is the branch of the git repository that contains the gitops config.
	// It cannot be empty.
	Branch string `yaml:"branch,omitempty"`

	// Path is the path to the directory or file that contains the gitops config.
	// It cannot be empty.
	Path string `yaml:"path,omitempty"`

	// Push specifies whether the gitops config is updated via git push.
	//
	// If false, gitimpart just clones the repository, may or may not update the gitops config locally,
	// and runs necessary commands to apply the changes (like kubectl-apply and terraform-apply).
	Push bool `yaml:"push,omitempty"`
}

type PullRequest struct{}

// RepositoryDispatch specifies whether the gitimpart run is triggered via GitHub repository_dispatch.
type RepositoryDispatch struct {
	// Owner is the owner of the repository that the repository_dispatch is sent to.
	Owner string `yaml:"owner"`
	// Repo is the name of the repository that the repository_dispatch is sent to.
	Repo string `yaml:"repo"`
}
