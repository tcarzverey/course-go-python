package urls

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
)

type ResponseCodeAggregator struct {
	client HttpClient
}

func NewURLAggregator(client HttpClient) *ResponseCodeAggregator {
	return &ResponseCodeAggregator{client: client}
}

func (u *ResponseCodeAggregator) Aggregate(ctx context.Context, urls <-chan string) (AggregationResult, error) {
	res := NewDefaultAggregationResult()
	if u.client == nil {
		u.client = http.DefaultClient
	}

	var wg sync.WaitGroup

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ustr, ok := <-urls:
				if !ok {
					return
				}
				wg.Add(1)
				go func(link string) {
					defer wg.Done()

					select {
					case <-ctx.Done():
						return
					default:
					}

					resp, err := u.client.Get(link)
					if err != nil {
						log.Printf("GET %s error: %v", link, err)
						return
					}
					io.Copy(io.Discard, resp.Body)
					_ = resp.Body.Close()

					res.add(resp.StatusCode)
				}(ustr)
			}
		}
	}()

	go func() {
		wg.Wait()
		res.markDone()
	}()

	return res, nil
}
