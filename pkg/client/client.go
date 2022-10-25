package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	baseURL    = "https://api.mapbox.com/geocoding/v5/mapbox.places"
	timeoutSec = 5
)

type Client interface {
	GetPostcode(longitude, latitude float64) (string, error)
}

type client struct {
	httpClient *http.Client
	token      string
}

type Resp struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Text string `json:"text"`
}

func New(token string) Client {
	return &client{
		httpClient: &http.Client{
			Timeout: timeoutSec * time.Second,
		},
		token: token,
	}
}

func (c *client) GetPostcode(longitude, latitude float64) (string, error) {
	resp, err := c.httpClient.Get(buildRequestURL(longitude, latitude, c.token))
	if err != nil {
		return "", fmt.Errorf("failed to send request to Mapbox: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get postcode from Mapbox: status code %d", resp.StatusCode)
	}
	simplifiedResp, err := getJSONResp(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to decode Mapbox response: %w", err)
	}
	postcode, err := simplifiedResp.GetText()
	if err != nil {
		return "", fmt.Errorf("unexpected Mapbox response: %w", err)
	}
	return postcode, nil
}

func buildRequestURL(longitude, latitude float64, token string) string {
	return baseURL + "/" +
		strconv.FormatFloat(longitude, 'f', 6, 64) + "," +
		strconv.FormatFloat(latitude, 'f', 6, 64) +
		".json?types=postcode&limit=1&access_token=" + token
}

func getJSONResp(body io.Reader) (*Resp, error) {
	simplifiedResp := &Resp{}
	if err := json.NewDecoder(body).Decode(simplifiedResp); err != nil {
		return nil, err
	}
	return simplifiedResp, nil
}

func (s Resp) GetText() (string, error) {
	if len(s.Features) == 0 {
		return "", errors.New("empty '.features' list")
	}
	if len(s.Features[0].Text) == 0 {
		return "", errors.New("empty '.features[0].text' field")
	}
	return s.Features[0].Text, nil
}
