package soze

import "fmt"

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
