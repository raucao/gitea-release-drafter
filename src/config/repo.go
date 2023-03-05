// this file implements configurations for repositories
package config

import (
	"io"

	"github.com/spf13/viper"
)

// RepoConfig holds all configurations for the repo we're running on
type RepoConfig struct {
	// DefaultBranch where we look for the configuration file
	DefaultBranch string `mapstructure:"default-branch"`
	// NameTemplate template for the release name
	NameTemplate string `mapstructure:"name-template"`
	// TagTemplate template for the release tag
	TagTemplate     string `mapstructure:"tag-template"`
	VersionResolver struct {
		Major struct {
			Labels []string
		}
		Minor struct {
			Labels []string
		}
		Patch struct {
			Labels []string
		}
		Default string
	} `mapstructure:"version-resolver"`
}

// ReadRepoConfig reads in the yaml config found in the default branch of the project and adds sensible defaults if values aren't set
func ReadRepoConfig(in io.Reader, defaultBranch string) (*RepoConfig, error) {
	vv := viper.New()
	vv.SetConfigType("yaml")

	err := vv.ReadConfig(in)
	if err != nil {
		return nil, err
	}

	// we set defaults here but if they are present in the configuration file they will be overwritten
	cfg := &RepoConfig{
		DefaultBranch: defaultBranch,
		NameTemplate:  "v$RESOLVED_VERSION",
		TagTemplate:   "v$RESOLVED_VERSION",
		VersionResolver: struct {
			Major   struct{ Labels []string }
			Minor   struct{ Labels []string }
			Patch   struct{ Labels []string }
			Default string
		}{
			Major:   struct{ Labels []string }{[]string{"major"}},
			Minor:   struct{ Labels []string }{[]string{"minor"}},
			Patch:   struct{ Labels []string }{[]string{"patch"}},
			Default: "minor",
		},
	}

	vv.Unmarshal(&cfg)

	return cfg, nil
}
