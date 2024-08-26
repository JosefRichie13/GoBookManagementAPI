package main

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// Generates an unique ID based on UUID.
// Strips out the hyphens and returns it as a string
func uniqueIDGenerator() string {

	id := uuid.New()
	return strings.Replace(id.String(), `-`, ``, -1)

}

// Performs sanitization on a string to mitigate XSS vulnerabilities
// Strips out < and > from a string and returns it
// Also strips out any extra spaces in the string
func sanitizeString(stringToSanitize string) string {

	replacingLeftArrow := strings.Replace(stringToSanitize, `<`, ``, -1)
	replacingRightArrow := strings.Replace(replacingLeftArrow, `>`, ``, -1)

	stripOutRegex := regexp.MustCompile(`\s+`)
	strippedOutString := stripOutRegex.ReplaceAllString(replacingRightArrow, " ")

	return strings.TrimSpace(strippedOutString)

}
