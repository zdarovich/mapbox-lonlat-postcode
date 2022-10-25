package model

type (
	Input struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lng"`
	}

	Output struct {
		Postcode string `json:"postcode"`
	}
)
