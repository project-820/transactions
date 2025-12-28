package platform

import (
	"fmt"

	"github.com/project-820/transactions/internal/platform/runner"
)

type Platform struct {
	runners []runner.Runner
}

func NewPlatform(runners ...runner.Runner) *Platform {
	return &Platform{
		runners: runners,
	}
}

func (s *Platform) Run() error {
	for _, runner := range s.runners {
		if err := runner.Start(); err != nil {
			return fmt.Errorf("start: %w", err)
		}
	}

	return nil
}

func (s *Platform) Stop() error {
	for _, runner := range s.runners {
		if err := runner.Stop(); err != nil {
			return fmt.Errorf("stop: %w", err)
		}
	}

	return nil
}
