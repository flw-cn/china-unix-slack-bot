package util

import (
	"regexp"
)

// NamedRegexpParse execute a regexp match and return all named group
func NamedRegexpParse(text string, exp *regexp.Regexp) (bool, map[string]string) {
	dict := make(map[string]string)

	allMatches := exp.FindStringSubmatch(text)
	if len(allMatches) == 0 {
		return false, dict
	}

	keys := exp.SubexpNames()
	if len(keys) != 0 {
		for i, name := range keys {
			if i != 0 && name != "" {
				dict[name] = allMatches[i]
			}
		}
	}

	return true, dict
}
