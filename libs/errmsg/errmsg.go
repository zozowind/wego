package errmsg

import "regexp"

var (
	enterRegexp1 = regexp.MustCompile(`\r\n`)
	enterRegexp2 = regexp.MustCompile(`\n`)
)

// ErrMsg struct
type ErrMsg struct {
	Code    int
	Message string
	Detail  string
}

// Get error message for implement error interface
func (e *ErrMsg) Error() string {
	return e.Message
}

func GetError(err *ErrMsg, detail string) *ErrMsg {
	var errCopy = *err
	errCopy.Detail = enterRegexp1.ReplaceAllString(detail, "")
	errCopy.Detail = enterRegexp2.ReplaceAllString(errCopy.Detail, "")
	return &errCopy
}
