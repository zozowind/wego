package errmsg

import (
	"fmt"
	"regexp"
)

var (
	enterRegexp1 = regexp.MustCompile(`\r\n`)
	enterRegexp2 = regexp.MustCompile(`\n`)
)

// ErrMsg struct
type ErrMsg struct {
	Code    int
	Message string
}

// Get error message for implement error interface
func (e *ErrMsg) Error() string {
	return e.Message
}

//GetError get common error message
func GetError(err *ErrMsg, detail string) error {
	detail = enterRegexp1.ReplaceAllString(detail, "")
	detail = enterRegexp2.ReplaceAllString(detail, "")
	return fmt.Errorf("[%d %s]%s", err.Code, err.Message, detail)
}
