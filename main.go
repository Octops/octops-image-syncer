package main

import (
	"context"
	"flag"
	"github.com/Octops/octops-image-syncer/cmd"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/signal"
	"time"
)

var (
	masterURL          string
	kubeconfig         string
	port               int
	syncPeriod         string
	metricsBindAddress string
)

func main() {
	flag.Parse()

	if kubeconfig == "" {
		kubeconfig = flag.Lookup("kubeconfig").Value.String()
	}

	clientConf, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		logrus.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	duration, err := time.ParseDuration(syncPeriod)
	if err != nil {
		logrus.WithError(err).Fatalf("error parsing sync-period flag: %s", syncPeriod)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := cmd.Execute(ctx, clientConf, duration, port, metricsBindAddress); err != nil {
		logrus.WithError(err).Fatal("failed to start syncer")
	}
}

func init() {
	if flag.Lookup("kubeconfig") == nil {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	}

	if flag.Lookup("master") == nil {
		flag.StringVar(&masterURL, "master", "", "The addr of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	}

	flag.IntVar(&port, "addr", 8090, "Port used by the broadcaster to communicate via http")
	flag.StringVar(&syncPeriod, "sync-period", "15s", "Determines the minimum frequency that the syncer will check for Fleets updates")
	flag.StringVar(&metricsBindAddress, "metrics-bind-address", "0.0.0.0:8095", "The TCP address that the controller should bind to for serving prometheus metrics")
}
