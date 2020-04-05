package entities

// Payload ...
type Payload struct {
	ID      string `json:"id"`
	Source  string `json:"source"`
	Subject string `json:"subject"`
	Body    []byte `json:"body"`
}
