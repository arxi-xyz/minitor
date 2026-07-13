package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"minitor/transport"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := transport.NewServer().Run(ctx); err != nil {
		log.Fatal(err)
	}
}
