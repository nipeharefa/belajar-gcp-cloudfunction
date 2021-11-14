// Package p contains an HTTP Cloud Function.
package p

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	Result struct {
		Code          uint16        `json:"code"`
		Latency       time.Duration `json:"latency"`
		LatencyString string        `json:"latency_string"`
		Timestamp     time.Time     `json:"timestamp"`
		BytesOut      uint64        `json:"bytes_out"`
		BytesIn       uint64        `json:"bytes_in"`
		RespHeader    http.Header   `json:"resp_header"`
	}
)

// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func HelloWorld(w http.ResponseWriter, r *http.Request) {

	var req *http.Request
	var err error
	var res *Result = new(Result)

	began := time.Now()

	client := &http.Client{}

	req, err = http.NewRequestWithContext(context.Background(), "GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	res.Timestamp = began.Add(time.Since(began))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, err)
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

	json.NewEncoder(w).Encode(res)
}
