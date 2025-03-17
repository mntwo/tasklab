package json

type Payload struct {
	Project     string            `json:"project"`
	Event       string            `json:"event"`
	Properties  map[string]string `json:"properties"`
	Type        string            `json:"type"`
	PassThrough []byte            `json:"pass_through"`
	Ts          int64             `json:"ts"`
}
