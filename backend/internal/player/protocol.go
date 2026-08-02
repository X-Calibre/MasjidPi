package player

type Command struct {
	Command []any `json:"command"`
}

type Response struct {
	Event string `json:"event,omitempty"`
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}
