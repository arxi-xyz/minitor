package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"minitor/config"
	"minitor/transport"
	"minitor/version"
)

func main() {
	configPath := flag.String("config", "", "path to JSON config file")
	addr := flag.String("addr", "", "server listen address")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version.String())
		return
	}

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
