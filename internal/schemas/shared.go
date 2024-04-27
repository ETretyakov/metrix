package schemas

type StatusMessage struct {
	Status bool   `json:"status"`
	Msg    string `json:"message"`
}

type WidgetResponse struct {
	Value uint64 `json:"value"`
}
