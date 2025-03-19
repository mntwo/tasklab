package dispatcher

import (
	"context"
	"errors"

	"github.com/mntwo/tasklab/encoding"
	"github.com/mntwo/tasklab/event_manager"
)

var (
	ErrEventManagerNotFound = errors.New("event manager not found")
)

func Dispatch(ctx context.Context, payload encoding.Payload) error {
	m, ok := event_manager.GetEventManager(payload.GetEvent())
	if !ok {
		return ErrEventManagerNotFound
	}
	m.Notify(payload.GetProperties())
	return nil
}
