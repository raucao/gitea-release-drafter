package src

import (
	"code.gitea.io/sdk/gitea"
)

func GetRepo(c *gitea.Client, owner string, repoName string) (*gitea.Repository, error) {
	repo, _, err := c.GetRepo(owner, repoName)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func FindReleases(c *gitea.Client, owner string, repo string) (*gitea.Release, *gitea.Release, error) {
	releases, _, err := c.ListReleases(owner, repo, gitea.ListReleasesOptions{})
	if err != nil {
		return nil, nil, err
	}

	var mostRecentRelease *gitea.Release
	var mostRecentDraftRelease *gitea.Release

	for _, r := range releases {
		if !r.IsPrerelease { // we don't care for pre-releases atm
			if r.IsDraft {
				if mostRecentDraftRelease == nil || r.CreatedAt.After(mostRecentDraftRelease.CreatedAt) {
					mostRecentDraftRelease = r
				}
			} else {
				if mostRecentRelease == nil || r.CreatedAt.After(mostRecentRelease.CreatedAt) {
					mostRecentRelease = r
				}
			}
		}
	}

	return mostRecentDraftRelease, mostRecentRelease, err
}

func CreateDraftRelease(c *gitea.Client, owner string, repo string, targetBranch string, version string, body string) (*gitea.Release, error) {
	release, _, err := c.CreateRelease(owner, repo, gitea.CreateReleaseOption{
		TagName:      version,
		Target:       targetBranch,
		Title:        version,
		Note:         body,
		IsDraft:      true,
		IsPrerelease: false,
	})
	if err != nil {
		return nil, err
	}

	return release, err
}

func UpdateExistingDraft(c *gitea.Client, owner string, repo string, draft *gitea.Release, nextVersion string, body string) (*gitea.Release, error) {
	rel, _, err := c.EditRelease(owner, repo, draft.ID, gitea.EditReleaseOption{
		TagName: nextVersion,
		Title:   nextVersion,
		Note:    body,
	})
	if err != nil {
		return nil, err
	}

	return rel, nil
}
