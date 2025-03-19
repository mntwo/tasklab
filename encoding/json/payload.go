package json

import (
	"encoding/json"

	"github.com/mntwo/tasklab/encoding"
)

var _ encoding.Payload = (*payload)(nil)

type payload struct {
	Project     string            `json:"project"`
	Event       string            `json:"event"`
	Properties  map[string]string `json:"properties"`
	Type        string            `json:"type"`
	PassThrough []byte            `json:"pass_through"`
	Ts          int64             `json:"ts"`
}

func New() encoding.Payload {
	return &payload{}
}

func (p *payload) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *payload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *payload) GetProject() string {
	return p.Project
}

func (p *payload) GetEvent() string {
	return p.Event
}

func (p *payload) GetProperties() map[string]string {
	return p.Properties
}

func (p *payload) GetType() string {
	return p.Type
}

func (p *payload) GetPassThrough() []byte {
	return p.PassThrough
}

func (p *payload) GetTs() int64 {
	return p.Ts
}
