package pool_test

import (
	"errors"
	"mapbox-lonlat-postcode/internal/model"
	"mapbox-lonlat-postcode/internal/pool"
	"mapbox-lonlat-postcode/pkg/logger"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockedMapboxClient struct {
	mock.Mock
}

func (m *MockedMapboxClient) GetPostcode(longitude, latitude float64) (string, error) {
	args := m.Called(longitude, latitude)
	return args.String(0), args.Error(1)
}

func prepareChannels(inSize, outSize int) (chan model.Input, chan model.Output) {
	input := make(chan model.Input, inSize)
	for i := 0; i < inSize; i++ {
		input <- model.Input{
			Longitude: float64(i),
			Latitude:  float64(i),
		}
	}
	output := make(chan model.Output, outSize)
	return input, output
}

func TestPool_Run_Succesful_Request_Must_Return_Postcode(t *testing.T) {
	bufSize := 2
	mockedClient := new(MockedMapboxClient)
	mockedClient.On("GetPostcode", mock.Anything, mock.Anything).Return("test", nil)

	input, _ := prepareChannels(bufSize, bufSize)
	l := logger.New("test", "debug")
	errHandler := func(err error) {
		l.Error("received error", err.Error())
	}
	p := pool.New(l, errHandler, mockedClient, bufSize)
	output := p.Run(input)

	for i := 0; i < bufSize; i++ {
		require.Equal(t, model.Output{
			Postcode: "test",
		}, <-output)
	}
	close(input)
	p.Wait()
}

func TestPool_Run_Failed_Request_Must_Return_Error(t *testing.T) {
	bufSize := 2
	mockedClient := new(MockedMapboxClient)
	mockedClient.On("GetPostcode", mock.Anything, mock.Anything).Return("", errors.New("test"))

	input, _ := prepareChannels(bufSize, bufSize)
	l := logger.New("test", "debug")

	errHandler := func(err error) {
		require.Equal(t, errors.New("test"), err)

	}
	p := pool.New(l, errHandler, mockedClient, bufSize)
	_ = p.Run(input)

	close(input)
	p.Wait()
}
