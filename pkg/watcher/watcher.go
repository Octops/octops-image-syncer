package watcher

import (
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"github.com/Octops/agones-event-broadcaster/pkg/broadcaster"
	"github.com/Octops/agones-event-broadcaster/pkg/brokers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"time"
)

type Config struct {
	ClientConfig   *rest.Config
	Duration       time.Duration
	Port           int
	MetricsAddress string
}

type ImageSyncer interface {
	brokers.Broker
}

type FleetWatcher struct {
	broadcaster *broadcaster.Broadcaster
}

func NewFleetWatcher(config *Config, imageSyncer ImageSyncer) (*FleetWatcher, error) {
	bc := broadcaster.New(config.ClientConfig, imageSyncer, config.Duration, config.Port, config.MetricsAddress)
	if err := bc.WithWatcherFor(&v1.Fleet{}).Build(); err != nil {
		return nil, errors.Wrap(err, "error creating broadcaster")
	}

	return &FleetWatcher{
		broadcaster: bc,
	}, nil
}

func (f *FleetWatcher) Start(ctx context.Context) error {
	go func() {
		if err := f.broadcaster.Start(); err != nil {
			logrus.WithError(err).Fatal("error starting broadcaster")
		}
	}()

	//TODO: refactor broadcaster to accept ctx on Start method
	<-ctx.Done()

	logrus.Info("shutting down syncer")
	return nil
}
