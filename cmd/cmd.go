package cmd

import (
	"context"
	"github.com/Octops/octops-image-syncer/pkg/syncer"
	"k8s.io/client-go/rest"
	"time"
)

func Execute(ctx context.Context, config *rest.Config, duration time.Duration, port int, metricsBindAddress string) error {
	syncer, err := syncer.NewFleetImageSyncer(config, duration, port, metricsBindAddress)
	if err != nil {
		return err
	}

	if err := syncer.Start(ctx); err != nil {
		return err
	}

	return nil
}
