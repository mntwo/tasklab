package application

import "errors"

var ErrApplicationClosed = errors.New("application closed")

type Application interface {
	Start() error
	Stop() error
	GetName() string
}
