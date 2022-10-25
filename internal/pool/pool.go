package pool

import (
	"mapbox-lonlat-postcode/internal/model"
	"mapbox-lonlat-postcode/pkg/client"
	"mapbox-lonlat-postcode/pkg/logger"
	"sync"
)

type Pool struct {
	workersCount int
	wg           *sync.WaitGroup
	mapboxClient client.Client
	output       chan model.Output
	logger       logger.Interface
	errHandler   func(err error)
}

func New(l logger.Interface, errorHandler func(err error), mapboxClient client.Client, workersCount int) *Pool {
	return &Pool{
		wg:           &sync.WaitGroup{},
		workersCount: workersCount,
		mapboxClient: mapboxClient,
		output:       make(chan model.Output),
		logger:       l,
		errHandler:   errorHandler,
	}
}

func (p *Pool) Run(in <-chan model.Input) <-chan model.Output {
	for i := 0; i < p.workersCount; i++ {
		p.wg.Add(1)
		go p.runTask(in)
	}
	return p.output
}

func (p *Pool) Wait() {
	p.wg.Wait()
	close(p.output)

}

func (p *Pool) runTask(input <-chan model.Input) {
	defer p.wg.Done()
	defer p.logger.Info("worker finished")
	for {
		in, ok := <-input
		if !ok {
			return
		}
		p.logger.Info("worker received input")
		postcode, err := p.mapboxClient.GetPostcode(in.Longitude, in.Latitude)
		if err != nil {
			p.errHandler(err)
			continue
		}
		out := model.Output{
			Postcode: postcode,
		}
		p.logger.Info("received response from mapbox")
		p.output <- out
	}
}
