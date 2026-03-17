package chat

type Message struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}
