package app

import (
	"context"
	"io"
	"mapbox-lonlat-postcode/internal/pool"
	"mapbox-lonlat-postcode/internal/stdio"
	"mapbox-lonlat-postcode/internal/util"
	"mapbox-lonlat-postcode/pkg/client"
	"mapbox-lonlat-postcode/pkg/logger"
)

type app struct {
	l            logger.Interface
	mapboxClient client.Client
	workersCount int
	reader       io.Reader
	writer       io.Writer
}

func New(l logger.Interface, mbc client.Client, wc int, r io.Reader, w io.Writer) *app {
	return &app{l, mbc, wc, r, w}
}

func (app *app) Run(ctx context.Context) {
	log := app.l
	ctx, cancel := context.WithCancel(ctx)

	errHandler := func(err error) {
		log.Error("received error", err.Error())
		cancel()
	}

	stdioIn := stdio.Reader(ctx, log, app.reader)
	in := util.BytesToInput(errHandler, stdioIn)

	workerPool := pool.New(log, errHandler, app.mapboxClient, app.workersCount)
	workerPoolOut := workerPool.Run(in)

	out := util.OutputToString(errHandler, workerPoolOut)

	done := stdio.Writer(errHandler, app.writer, out)

	workerPool.Wait()

	<-done
	log.Info("job finished")
}
