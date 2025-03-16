package handler

import (
	"context"

	"github.com/mntwo/tasklab/gen_event"
	"github.com/mntwo/tasklab/internal/log"
	"go.uber.org/zap"
)

var _ gen_event.Handler = (*SampleA)(nil)

type SampleA struct{}

func (s *SampleA) Init() {
	log.Info(context.Background(), "SampleA init")
}

func (s *SampleA) HandleEvent(ctx context.Context, event gen_event.Event) {
	// Add your logic here
	log.Info(ctx, "SampleA handle event", zap.Any("event", event))
}

func (s *SampleA) Close() error {
	log.Info(context.Background(), "SampleA close")
	return nil
}

var _ gen_event.Handler = (*SampleB)(nil)

type SampleB struct{}

func (s *SampleB) Init() {
	log.Info(context.Background(), "SampleB init")
}

func (s *SampleB) HandleEvent(ctx context.Context, event gen_event.Event) {
	// Add your
	log.Info(ctx, "SampleB handle event", zap.Any("event", event))
}

func (s *SampleB) Close() error {
	log.Info(context.Background(), "SampleB close")
	return nil
}
