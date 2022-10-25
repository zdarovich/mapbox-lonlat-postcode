package util

import (
	"bytes"
	"encoding/json"
	"mapbox-lonlat-postcode/internal/model"
)

func OutputToString(errHandler func(err error), in <-chan model.Output) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			output, ok := <-in
			if !ok {
				return
			}
			raw, err := json.Marshal(&output)
			if err != nil {
				errHandler(err)
				continue
			}
			compact := new(bytes.Buffer)
			if err := json.Compact(compact, raw); err != nil {
				errHandler(err)
				continue
			}
			out <- compact.String()
		}
	}()
	return out
}

func BytesToInput(errHandler func(err error), in <-chan []byte) <-chan model.Input {
	out := make(chan model.Input)
	go func() {
		defer close(out)
		for {
			output, ok := <-in
			if !ok {
				return
			}
			input := model.Input{}
			if err := json.Unmarshal(output, &input); err != nil {
				errHandler(err)
				continue
			}
			out <- input
		}
	}()
	return out
}
