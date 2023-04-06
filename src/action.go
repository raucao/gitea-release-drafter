package src

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
	"git.andinfinity.de/gitea-release-drafter/src/config"
	"github.com/sethvargo/go-githubactions"
)

type Action struct {
	config *config.DrafterConfig

	globalContext context.Context
	client        *gitea.Client
}

// NewAction factory for a new action
func NewAction(ctx *context.Context, cfg *config.DrafterConfig) (*Action, error) {
	gitea, err := gitea.NewClient(cfg.ApiUrl, gitea.SetToken(cfg.Token))
	if err != nil {
		return nil, err
	}

	return &Action{
		config:        cfg,
		globalContext: *ctx,
		client:        gitea,
	}, nil
}

// updateOrCreateDraftRelease
func updateOrCreateDraftRelease(a *Action, cfg *config.RepoConfig) (*gitea.Release, error) {
	draft, last, err := FindReleases(a.client, a.config.RepoOwner, a.config.RepoName)
	if err != nil {
		return nil, err
	}

	changelog, err := GenerateChangelog(a.client, a.config.RepoOwner, a.config.RepoName, last)
	if err != nil {
		return nil, err
	}

	if len(*changelog) == 0 {
		githubactions.Infof("No updates found")
		return nil, nil
	}

	// render changelog
	var b strings.Builder

	b.WriteString("# What's Changed")
	b.WriteString("\n\n")

	// TODO: group by given label categories in config
	// default to dumping everything by date desc.

	if changelog != nil {
		for label, prs := range *changelog {
			if len(prs) > 0 {
				// TODO: here we should take the label from the config and only default to the name
				fmt.Fprintf(&b, "## %s\n\n", strings.Title(label))

				for _, pr := range prs {
					fmt.Fprintf(&b, "* %s (#%d) @%s\n", pr.Title, pr.Index, pr.Poster.UserName)
				}

				b.WriteString("\n")
			}
		}
	}

	nextVersion, err := ResolveVersion(cfg, last, changelog)
	if err != nil {
		return nil, err
	}

	title := FillVariables(cfg.NameTemplate, TemplateVariables{
		ReleaseVersion: nextVersion.String(),
	})

	// FIXME: require RESOLVED_VERSION to be set?
	tag := FillVariables(cfg.TagTemplate, TemplateVariables{
		ReleaseVersion: nextVersion.String(),
	})

	if draft != nil {
		updatedDraft, err := UpdateExistingDraft(a.client, a.config.RepoOwner, a.config.RepoName, draft, title, tag, b.String())
		if err != nil {
			return nil, err
		}

		return updatedDraft, nil
	}

	newDraft, err := CreateDraftRelease(a.client, a.config.RepoOwner, a.config.RepoName, cfg.DefaultBranch, title, tag, b.String())
	if err != nil {
		return nil, err
	}

	return newDraft, nil
}

// GetConfigFile reads the local configuration file in `.gitea/` in the ref branch (probably main/master)
func (a *Action) GetConfigFile(ref string) (*bytes.Reader, error) {
	data, _, err := a.client.GetFile(a.config.RepoOwner, a.config.RepoName, ref, a.config.ConfigPath)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), err
}

// Run builds the configuration and executes the action logic
func (a *Action) Run() error {
	// fetch the repo to retrieve the default branch to be set as the config default
	repo, err := GetRepo(a.client, a.config.RepoOwner, a.config.RepoName)
	if err != nil {
		return err
	}

	githubactions.Debugf("Found default branch %s", repo.DefaultBranch)

	// build repo config
	configReader, err := a.GetConfigFile(repo.DefaultBranch)
	if err != nil {
		if err.Error() != "404 Not Found" {
			return err
		} else {
			// no config file found
			githubactions.Warningf("No such config file: .gitea/%s", a.config.ConfigPath)
			configReader = bytes.NewReader([]byte{})
		}
	}

	config, err := config.ReadRepoConfig(configReader, repo.DefaultBranch)
	if err != nil {
		return err
	}

	draft, err := updateOrCreateDraftRelease(a, config)
	if err != nil {
		return err
	}

	if draft != nil {
		githubactions.Infof("created draft release %s", draft.Title)
	}

	return nil
}
