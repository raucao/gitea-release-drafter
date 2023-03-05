package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidConfigShouldWork(t *testing.T) {
	// Given
	// A valid yml config
	config := `
name-template: 'v$RESOLVED_VERSION 🌈'
tag-template: 'tag-v$RESOLVED_VERSION'
version-resolver:
  major:
    labels:
      - 'major-test'
  minor:
    labels:
      - 'minor-test'
  patch:
    labels:
      - 'patch-test'
  default: 'minor-test'`

	in := strings.NewReader(config)

	// When
	// Reading in the config
	cfg, err := ReadRepoConfig(in, "main")

	// Then
	// No error should've occurred
	assert.NoError(t, err)

	// The name template should've been read in properly
	assert.Equal(t, "v$RESOLVED_VERSION 🌈", cfg.NameTemplate)

	// The name template should've been read in properly
	assert.Equal(t, "tag-v$RESOLVED_VERSION", cfg.TagTemplate)

	// The version resolver major labels should've been read in properly
	assert.Len(t, cfg.VersionResolver.Major.Labels, 1)
	assert.Contains(t, cfg.VersionResolver.Major.Labels, "major-test")

	// The version resolver minor labels should've been read in properly
	assert.Len(t, cfg.VersionResolver.Minor.Labels, 1)
	assert.Contains(t, cfg.VersionResolver.Minor.Labels, "minor-test")

	// The version resolver patch labels should've been read in properly
	assert.Len(t, cfg.VersionResolver.Patch.Labels, 1)
	assert.Contains(t, cfg.VersionResolver.Patch.Labels, "patch-test")

	// The version resolver default should've been read in properly
	assert.Equal(t, "minor-test", cfg.VersionResolver.Default)
}

func TestEmptyConfigShouldUseDefaults(t *testing.T) {
	// Given
	// An empty yml config
	config := ``

	in := strings.NewReader(config)

	// When
	// Reading in the config
	cfg, err := ReadRepoConfig(in, "main")

	// Then
	// No error should've occurred
	assert.NoError(t, err)

	// The name template should've been read in properly
	assert.Equal(t, "v$RESOLVED_VERSION", cfg.NameTemplate)

	// The name template should've been read in properly
	assert.Equal(t, "v$RESOLVED_VERSION", cfg.TagTemplate)

	// The version resolver major labels should've been read in properly
	assert.Len(t, cfg.VersionResolver.Major.Labels, 1)
	assert.Contains(t, cfg.VersionResolver.Major.Labels, "major")

	// The version resolver minor labels should've been read in properly
	assert.Len(t, cfg.VersionResolver.Minor.Labels, 1)
	assert.Contains(t, cfg.VersionResolver.Minor.Labels, "minor")

	// The version resolver patch labels should've been read in properly
	assert.Len(t, cfg.VersionResolver.Patch.Labels, 1)
	assert.Contains(t, cfg.VersionResolver.Patch.Labels, "patch")

	// The version resolver default should've been read in properly
	assert.Equal(t, "minor", cfg.VersionResolver.Default)
}
