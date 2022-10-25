package app_test

import (
	"bytes"
	"context"
	"mapbox-lonlat-postcode/internal/app"
	"mapbox-lonlat-postcode/pkg/logger"
	"strings"
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

func TestRun_Valid_Input_Must_Output_Successfully(t *testing.T) {
	inputData := strings.Join([]string{
		`{"lat": 0.1, "lng": 0.1}`,
		`{"lat": 0.2, "lng": 0.2}`,
	}, "\n")
	outputData := strings.Join([]string{
		`{"postcode":"test"}`,
		`{"postcode":"test"}`,
	}, "\n") + "\n"
	mbc := new(MockedMapboxClient)
	mbc.On("GetPostcode", mock.Anything, mock.Anything).Return("test", nil)
	outBuffer := new(bytes.Buffer)
	inputReader := strings.NewReader(inputData)
	testApp := app.New(logger.New("test", "debug"), mbc, 3, inputReader, outBuffer)
	testApp.Run(context.Background())
	require.Equal(t, outputData, outBuffer.String())
}

func TestRun_Invalid_Input_Must_Cancel_Gracefully_And_Receive_1_Valid_Result(t *testing.T) {
	inputData := strings.Join([]string{
		`{"lat": 0.1, "lng": 0.1}`,
		`{"lat": 0.2,`,
	}, "\n")
	outputData := strings.Join([]string{
		`{"postcode":"test"}`,
	}, "\n") + "\n"
	mbc := new(MockedMapboxClient)
	mbc.On("GetPostcode", mock.Anything, mock.Anything).Return("test", nil)
	outBuffer := new(bytes.Buffer)
	inputReader := strings.NewReader(inputData)
	testApp := app.New(logger.New("test", "debug"), mbc, 3, inputReader, outBuffer)
	testApp.Run(context.Background())
	require.Equal(t, outputData, outBuffer.String())
}
