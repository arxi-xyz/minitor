package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"minitor/config"
	"minitor/transport"
)

func main() {
	configPath := flag.String("config", "", "path to JSON config file")
	addr := flag.String("addr", "", "server listen address")
	flag.Parse()

	cfg, err := config.Load(config.LoadOptions{
		Path:       *configPath,
		ServerAddr: *addr,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := transport.NewServer(cfg).Run(ctx); err != nil {
		log.Fatal(err)
	}
}
