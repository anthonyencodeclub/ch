package errfmt

import "fmt"

// Format returns a user-friendly error string.
func Format(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("Error: %s", err.Error())
}
