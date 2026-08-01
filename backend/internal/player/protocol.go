package player

type Command struct {
	Command []any `json:"command"`
}

type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}
