package encoding

type Payload interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
	GetProject() string
	GetEvent() string
	GetProperties() map[string]string
	GetType() string
	GetPassThrough() []byte
	GetTs() int64
}
