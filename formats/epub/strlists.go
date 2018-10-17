package epub

import (
	"strings"
)

// ExistsStr checks whether a string str contained in list (case sensitive)
func listContainsString(str string, list []string) bool {

	if len(list) == 0 {
		return false
	}

	for _, item := range list {
		if item == str {
			return true
		}
	}

	return false
}

// ExistsText checks whether a string str contained in list (case insensitive)
func listContainsText(str string, list []string) bool {

	if len(list) == 0 {
		return false
	}

	for _, item := range list {
		if strings.ToLower(item) == strings.ToLower(str) {
			return true
		}
	}

	return false
}
