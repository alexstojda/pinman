package utils

// Source: https://github.com/metal3d/go-slugify/blob/master/main.go
// Package slugify provide a function that
// gives a non accentuated and minus separated string from a
// accentuated string. The code is based from a Javascript function
// that you can get here:
// http://irz.fr/slugme-permalien-javascript-slug/

import (
	"regexp"
	"strings"
)

// Replacement structure
type replacement struct {
	re *regexp.Regexp
	ch string
}

// Build regexps and replacements
var (
	rExps = []replacement{
		{re: regexp.MustCompile(`[\xC0-\xC6]`), ch: "A"},
		{re: regexp.MustCompile(`[\xE0-\xE6]`), ch: "a"},
		{re: regexp.MustCompile(`[\xC8-\xCB]`), ch: "E"},
		{re: regexp.MustCompile(`[\xE8-\xEB]`), ch: "e"},
		{re: regexp.MustCompile(`[\xCC-\xCF]`), ch: "I"},
		{re: regexp.MustCompile(`[\xEC-\xEF]`), ch: "i"},
		{re: regexp.MustCompile(`[\xD2-\xD6]`), ch: "O"},
		{re: regexp.MustCompile(`[\xF2-\xF6]`), ch: "o"},
		{re: regexp.MustCompile(`[\xD9-\xDC]`), ch: "U"},
		{re: regexp.MustCompile(`[\xF9-\xFC]`), ch: "u"},
		{re: regexp.MustCompile(`[\xC7-\xE7]`), ch: "c"},
		{re: regexp.MustCompile(`\xD1`), ch: "N"},
		{re: regexp.MustCompile(`\xF1`), ch: "n"},
	}
	spacereg       = regexp.MustCompile(`\s+`)
	noncharreg     = regexp.MustCompile(`[^A-Za-z0-9-]`)
	minusrepeatreg = regexp.MustCompile(`-{2,}`)
)

// Slugify function returns slugifies string "s"
func Slugify(s string, maxLength ...int) string {
	for _, r := range rExps {
		s = r.re.ReplaceAllString(s, r.ch)
	}

	s = strings.ToLower(s)
	s = spacereg.ReplaceAllString(s, "-")
	s = noncharreg.ReplaceAllString(s, "")
	s = minusrepeatreg.ReplaceAllString(s, "-")

	if len(maxLength) > 0 && maxLength[0] > 0 && len(s) > maxLength[0] {
		s = s[:maxLength[0]]
	}

	return s
}
