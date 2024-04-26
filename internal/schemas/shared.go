package schemas

type StatusMessage struct {
	Status bool   `json:"status"`
	Msg    string `json:"message"`
}
