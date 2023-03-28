package src

import (
	"code.gitea.io/sdk/gitea"
	"git.andinfinity.de/gitea-release-drafter/src/config"
	"github.com/Masterminds/semver"
	"golang.org/x/exp/slices"
)

// ResolveVersion determines the next version to be used for a release depending on the labels used in the pull requests merged after the last release.
func ResolveVersion(cfg *config.RepoConfig, last *gitea.Release, changelog *Changelog) (*semver.Version, error) {
	// determine next version
	var nextVersion semver.Version

	// no prior release, starting with "v0.1.0"
	if last == nil {
		ver, err := semver.NewVersion("0.1")
		if err != nil {
			return nil, err
		}
		nextVersion = *ver
	} else {
		lastVersion, err := semver.NewVersion(last.TagName) // FIXME: what if it's not the version?
		if err != nil {
			return nil, err
		}

		incMajor := false
		incMinor := false
		incPatch := false

		// check labels
		for k := range *changelog {
			if slices.Contains(cfg.VersionResolver.Major.Labels, k) {
				incMajor = true
			}
			if slices.Contains(cfg.VersionResolver.Minor.Labels, k) {
				incMinor = true
			}
			if slices.Contains(cfg.VersionResolver.Patch.Labels, k) {
				incPatch = true
			}
		}

		if incMajor {
			nextVersion = lastVersion.IncMajor()
		} else if incMinor {
			nextVersion = lastVersion.IncMinor()
		} else if incPatch {
			nextVersion = lastVersion.IncPatch()
		} else {
			// default
			if cfg.VersionResolver.Default == "major" {
				nextVersion = lastVersion.IncMajor()
			} else if cfg.VersionResolver.Default == "minor" {
				nextVersion = lastVersion.IncMinor()
			} else {
				nextVersion = lastVersion.IncPatch()
			}
		}
	}

	return &nextVersion, nil
}
