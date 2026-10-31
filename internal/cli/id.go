package cli

import "strconv"

// parsePositiveID accepts only decimal positive integers for use in API path segments.
func parsePositiveID(s string) (string, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil || n == 0 {
		return "", usageErr("ID must be a positive integer")
	}
	return strconv.FormatUint(n, 10), nil
}
