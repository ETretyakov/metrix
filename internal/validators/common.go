// Module "validators" aggregates all necessary validations for web-service handlers.
package validators

import "fmt"

// ParsingValueError - the structure for value error for parsing.
type ParsingValueError struct {
	msg string
}

// Error - the method to match error interface.
func (v ParsingValueError) Error() string {
	return v.msg
}

// NewParsingValueError - the builder function for ParsingValueError.
func NewParsingValueError(msg string, vars ...any) ParsingValueError {
	return ParsingValueError{
		msg: fmt.Sprintf(msg, vars...),
	}
}
