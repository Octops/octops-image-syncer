package cmd

import (
	"context"
	"github.com/Octops/octops-image-syncer/pkg/syncer"
	"github.com/Octops/octops-image-syncer/pkg/transport"
	"github.com/Octops/octops-image-syncer/pkg/watcher"
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
	"os"
	"time"
)

func Execute(ctx context.Context, config *rest.Config, duration time.Duration, port int, metricsBindAddress string) error {
	target := os.Getenv("CONN_TARGET")
	conn, err := transport.NewConn(target)
	if err != nil {
		return errors.Wrapf(err, "failed to create connection to: %s", target)
	}
	defer conn.Close()

	imageSyncer := syncer.NewFleetImageSyncer(conn)
	watcherConfig := &watcher.Config{
		ClientConfig:   config,
		Duration:       duration,
		Port:           port,
		MetricsAddress: metricsBindAddress,
	}

	fleetWatcher, err := watcher.NewFleetWatcher(watcherConfig, imageSyncer)
	if err != nil {
		return errors.Wrap(err, "failed to create watcher")
	}

	if err := fleetWatcher.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start fleet watcher")
	}

	return nil
}
