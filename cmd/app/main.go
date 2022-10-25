package main

import (
	"context"
	"flag"
	"mapbox-lonlat-postcode/internal/app"
	"mapbox-lonlat-postcode/pkg/client"
	"mapbox-lonlat-postcode/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := flag.String("token", "", "token")
	workerCount := flag.Int("workers", 1, "workers")
	flag.Parse()
	l := logger.New("prod", "debug")
	ctx, cancel := context.WithCancel(context.Background())

	interrupt := make(chan os.Signal, 1)
	defer close(interrupt)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case s := <-interrupt:
			l.Info("app - Run - signal: " + s.String())
			cancel()
		}
	}()
	mbc := client.New(*token)
	testApp := app.New(l, mbc, *workerCount, os.Stdin, os.Stdout)
	testApp.Run(ctx)
}
