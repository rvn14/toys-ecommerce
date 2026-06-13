package utils

import (
	"regexp"
	"strings"
)

var nonAlphaNumericRegex = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(value string) string {
	slug := strings.ToLower(strings.TrimSpace(value))
	slug = nonAlphaNumericRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	if slug == "" {
		return "item"
	}

	return slug
}
