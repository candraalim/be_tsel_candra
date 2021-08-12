package util

import (
	"regexp"
	"strings"
)

var regexMsisdn = regexp.MustCompile("^62[0-9]{7,14}$").MatchString

func ValidateAndSanitizeMsisdn(msisdn string) (string, error) {
	msisdn = strings.TrimSpace(msisdn)
	if strings.HasPrefix(msisdn, "08") {
		msisdn = "628" + strings.TrimPrefix(msisdn, "08")
	}
	if regexMsisdn(msisdn) {
		return msisdn, nil
	}
	return "", ErrorInvalidRequest
}
