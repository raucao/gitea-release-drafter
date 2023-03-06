package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceVariables(t *testing.T) {
	// Given
	// valid variable values
	v := TemplateVariables{
		ReleaseVersion: "1.2.3",
	}

	// and a string containing all variables we have
	str := "tag-v$RESOLVED_VERSION-foo"

	// When
	// filling the variables
	res := FillVariables(str, v)

	// Then
	// the string should've been filled as expected
	assert.Equal(t, "tag-v1.2.3-foo", res)
}

func TestReplaceVariablesShouldSkipIfVarsAreAbsent(t *testing.T) {
	// Given
	// valid variable values
	v := TemplateVariables{
		ReleaseVersion: "1.2.3",
	}

	// and a string containing no variables
	str := "tag-v-foo"

	// When
	// filling the variables
	res := FillVariables(str, v)

	// Then
	// the string should'nt have been filled
	assert.Equal(t, "tag-v-foo", res)
}
