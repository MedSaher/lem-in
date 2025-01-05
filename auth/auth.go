package auth

import (
	// "errors"
	"regexp"
)

// is the name valid:
func IsValidFile(fileName string) bool {
	pattern := `^(example0[0-8]|badexample0[0-1])\.txt$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(fileName)
}
