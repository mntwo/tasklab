package application

type Application interface {
	Start() error
	Stop() error
	GetName() string
}
