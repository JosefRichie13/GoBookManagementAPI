package main

import (
	"regexp"
	"strings"
	"time"

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

// Checks if a string is in DD-MMM-YYYY format
// Returns TRUE if yes, or FALSE if not
func checkDateFormat(date string) bool {

	regex := regexp.MustCompile(`^[0-9]{2}-[A-Za-z]{3}-[0-9]{4}$`)
	if !regex.MatchString(date) {
		return false
	}
	_, err := time.Parse("02-Jan-2006", date)

	return err == nil

}

// Converts date in DD-MMM-YYYY format to Epoch Time and returns it
func convertDateToEpoch(date string) int {

	layout := "02-Jan-2006"
	t, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}

	epoch := t.Unix()
	return int(epoch)

}

// Converts Epoch time to date in DD-MMM-YYYY format and returns it
func convertEpochToDate(epochTime *int) string {

	// If the received time is 0, return an empty string. If we do not have this, then time returned will be 01-Jan-1970
	if int64(*epochTime) == 0 {
		return ""
	}

	unixTime := int64(*epochTime)

	// Convert Unix time to Time object
	t := time.Unix(unixTime, 0)

	// Format the time to DD-MMM-YYYY
	return t.Format("02-Jan-2006")

}
