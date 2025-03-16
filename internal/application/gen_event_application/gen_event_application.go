package gen_event_application

import (
	"github.com/mntwo/tasklab/event_manager"
	_ "github.com/mntwo/tasklab/event_manager"
	"github.com/mntwo/tasklab/gen_event"
	"github.com/mntwo/tasklab/handler"
	"github.com/mntwo/tasklab/internal/application"
)

var _ application.Application = (*GenEventApplication)(nil)

type GenEventApplication struct {
	Name   string
	stopCh chan struct{}
}

func New(name string) *GenEventApplication {
	return &GenEventApplication{
		Name:   name,
		stopCh: make(chan struct{}),
	}
}

func (a *GenEventApplication) Start() error {
	// This is a sample task that will be added to the event manager when the application starts.
	event_manager.AddEventManager("sample_task", gen_event.NewEventManager(10))
	m, ok := event_manager.GetEventManager("sample_task")
	if ok {
		m.AddHandler(&handler.SampleA{})
		m.AddHandler(&handler.SampleB{})
	}
	<-a.stopCh
	return application.ErrApplicationClosed
}

func (a *GenEventApplication) Stop() error {
	event_manager.Stop()
	close(a.stopCh)
	return nil
}

func (a *GenEventApplication) GetName() string {
	return a.Name
}
