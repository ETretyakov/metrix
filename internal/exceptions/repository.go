package exceptions

type RecordNotFound struct {
	Msg string
}

func (e RecordNotFound) Error() string {
	return e.Msg
}
