package stdio

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mapbox-lonlat-postcode/pkg/logger"
)

func Reader(ctx context.Context, logger logger.Interface, r io.Reader) <-chan []byte {
	out := make(chan []byte)
	go func() {
		defer close(out)
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				logger.Info("context was cancelled")
				break
			default:
			}
			raw := scanner.Bytes()
			out <- raw
		}
		if err := scanner.Err(); err != nil {
			logger.Error("failed to read input stream: %s", err.Error())
		}
	}()
	return out
}

func Writer(errHandler func(err error), w io.Writer, in <-chan string) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case output, ok := <-in:
				if !ok {
					return
				}
				_, err := fmt.Fprintln(w, output)
				if err != nil {
					errHandler(err)
					continue
				}
			default:
			}

		}
	}()
	return done
}
