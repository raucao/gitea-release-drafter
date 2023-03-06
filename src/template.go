package src

import (
	"reflect"
	"strings"
)

const varTag = "var"

type TemplateVariables struct {
	ReleaseVersion string `var:"$RESOLVED_VERSION"`
}

func FillVariables(str string, vars TemplateVariables) string {
	t := reflect.TypeOf(vars)
	v := reflect.ValueOf(vars)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		variable := field.Tag.Get(varTag)
		varVal := v.Field(i).String()
		str = strings.ReplaceAll(str, variable, varVal)
	}

	return str
}
