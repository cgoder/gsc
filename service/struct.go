package service

type GscMsg struct {
	Flag string
	Msg  Message
}

// Message payload from client.
type Message struct {
	Type    string `json:"type"`
	Input   string `json:"input"`
	Output  string `json:"output"`
	Payload string `json:"payload"`
}

// Status response to client.
type Status struct {
	Progress int     `json:"progress"`
	Percent  float64 `json:"percent"`
	Speed    string  `json:"speed"`
	FPS      float64 `json:"fps"`
	Err      string  `json:"err,omitempty"`
}
