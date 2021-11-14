// Package p contains an HTTP Cloud Function.
package p

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type (
	ResultConcurrent struct {
		TotalRequest      uint16         `json:"total_request"`
		StatusCodeCounter map[uint16]int `json:"status_code_counter"`
		TargetURL         string         `json:"target_url"`
	}
	Result struct {
		Code          uint16        `json:"code"`
		Latency       time.Duration `json:"latency"`
		LatencyString string        `json:"latency_string"`
		Timestamp     time.Time     `json:"timestamp"`
		BytesOut      uint64        `json:"bytes_out"`
		BytesIn       uint64        `json:"bytes_in"`
		RespHeader    http.Header   `json:"resp_header"`
		Error         error         `json:"error"`
	}
)

const (
	targetURL string = "https://jsonplaceholder.typicode.com/posts/1"
)

func hit(wg *sync.WaitGroup, respChan chan *Result) {
	var req *http.Request
	var err error
	var res *Result = new(Result)

	defer wg.Done()

	began := time.Now()

	client := &http.Client{}

	req, err = http.NewRequestWithContext(context.Background(), "GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		res.Error = err
		respChan <- res
		return
	}

	res.Timestamp = began.Add(time.Since(began))
	resp, err := client.Do(req)
	if err != nil {
		res.Error = err
		respChan <- res
		return
	}

	res.Latency = time.Since(res.Timestamp)
	res.LatencyString = res.Latency.String()
	res.Code = uint16(resp.StatusCode)

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if req.ContentLength != -1 {
		res.BytesOut = uint64(req.ContentLength)
	}

	var buf bytes.Buffer

	io.Copy(&buf, resp.Body)
	res.BytesIn = uint64(buf.Len())

	respChan <- res

}

// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func HelloWorld(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup
	resp := new(ResultConcurrent)

	counter := 10

	counterFromReq := r.URL.Query().Get("counter")

	if counterFromReq != "" {
		counter, _ = strconv.Atoi(counterFromReq)
		if counter < 1 {
			counter = 1
		}
	}

	resultChan := make(chan *Result, counter)

	resp.StatusCodeCounter = make(map[uint16]int)
	resp.TotalRequest = uint16(counter)
	resp.TargetURL = targetURL

	wg.Add(counter)
	for i := 0; i < counter; i++ {
		go hit(&wg, resultChan)
	}

	wg.Wait()
	close(resultChan)

	for v := range resultChan {
		statusCode := v.Code

		if _, ok := resp.StatusCodeCounter[statusCode]; ok {
			resp.StatusCodeCounter[statusCode]++
		} else {
			resp.StatusCodeCounter[statusCode] = 1
		}
	}

	json.NewEncoder(w).Encode(resp)
}
