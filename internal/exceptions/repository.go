package exceptions

type RecordNotFoundError struct {
	Msg string
}

func (e RecordNotFoundError) Error() string {
	return e.Msg
}
