package app

import (
	"github.com/mntwo/tasklab/internal/application"
	"github.com/mntwo/tasklab/internal/application/gen_event_application"
)

func NewGenEventApp() application.Application {
	return gen_event_application.New("gen_event_app")
}
