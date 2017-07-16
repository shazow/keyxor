package soze

import "fmt"

// ErrSizeMismatch is returned when the written output or read input does not
// match the size we were expecting.
type ErrSizeMismatch struct {
	Message  string
	Actual   int
	Expected int
}

func (err ErrSizeMismatch) Error() string {
	if err.Actual == 0 && err.Expected == 0 {
		return fmt.Sprintf("%s: missing data", err.Message)
	}
	return fmt.Sprintf("%s: mismatched size: %d != %d", err.Message, err.Actual, err.Expected)
}
