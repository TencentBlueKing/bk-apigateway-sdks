package internal

import (
	"regexp"
	"strings"
)

var placeholderRe *regexp.Regexp

// ReplacePlaceHolder replaces the placeholder with the given string.
func ReplacePlaceHolder(s string, params map[string]string) string {
	return placeholderRe.ReplaceAllStringFunc(
		s, func(placeholder string) string {
			key := strings.Trim(placeholder, "{ }")
			value, ok := params[key]
			if !ok {
				return placeholder
			}

			return value
		},
	)
}

func init() {
	// available placeholder pattern: {param} or { param }
	placeholderRe = regexp.MustCompile(`{\s*.*?\s*}`)
}
