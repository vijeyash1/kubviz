package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/intelops/kubviz/agent/git/pkg/application"
	"github.com/intelops/kubviz/agent/git/pkg/clients"
	"github.com/intelops/kubviz/agent/git/pkg/config"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg := &config.Config{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Could not parse env Config: %v", err)
	}

	// Connect to NATS
	natsContext, err := clients.NewNATSContext(cfg)
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	app := application.New(cfg, natsContext)

	go app.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	app.Close()
}
