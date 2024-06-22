package validators

import "fmt"

type ParsingValueError struct {
	msg string
}

func (v ParsingValueError) Error() string {
	return v.msg
}

func NewParsingValueError(msg string, vars ...any) ParsingValueError {
	return ParsingValueError{
		msg: fmt.Sprintf(msg, vars...),
	}
}
