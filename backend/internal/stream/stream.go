package stream

type Stream struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	URL      string `json:"url"`
}
