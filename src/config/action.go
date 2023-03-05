package config

import (
	"os"

	githubactions "github.com/sethvargo/go-githubactions"
)

// DrafterConfig holds all configurations we need for the drafter to run
type DrafterConfig struct {
	// RepoOwner as provided by the github context
	RepoOwner string
	// RepoName as provided by the github context
	RepoName string
	// ApiUrl of gitea as provided by the "GITHUB_SERVER_URL" env var
	ApiUrl string
	// Token as provided by the "GITHUB_TOKEN" env var
	Token string
	// ConfigPath as provided by the "config-path" action input. defaults to ".gitea/release-drafter.yml"
	ConfigPath string
}

// NewFromInputs creates a new drafter config by using the action inputs and the github context
func NewFromInputs(action *githubactions.Action) (*DrafterConfig, error) {
	actionCtx, err := action.Context()
	if err != nil {
		return nil, err
	}

	var configPath string
	inputConfigPath := action.GetInput("config-path")

	if inputConfigPath == "" {
		configPath = ".gitea/release-drafter.yml"
	} else {
		configPath = inputConfigPath
	}

	owner, name := actionCtx.Repo()
	c := DrafterConfig{
		RepoOwner:  owner,
		RepoName:   name,
		ApiUrl:     actionCtx.ServerURL,
		Token:      os.Getenv("GITHUB_TOKEN"),
		ConfigPath: configPath,
	}

	return &c, nil
}
