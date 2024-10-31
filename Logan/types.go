package Logan

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	Reset  = "\033[0m"
	Green  = "\033[32m" //app
	Blue   = "\033[34m" // server
	Cyan   = "\033[36m" // HTTP
	Yellow = "\033[33m" // nil
)

type Logan interface {
	AppSide(msg string)
	ServerSide(msg string)
	HTTPPrint(r *http.Request, resp *http.Response, target string)
}

type Log struct {
	Category string
}

func (l *Log) AppSide(msg string) {
	fmt.Printf("%s[%s] [%s] [APP] %s%s\n", Green, time.Now().Format("15:04:05"), l.Category, msg, Reset)
}

func (l *Log) ServerSide(msg string) {
	fmt.Printf("%s[%s] [%s] [SERVER] %s%s\n", Blue, time.Now().Format("15:04:05"), l.Category, msg, Reset)
}

func (l *Log) HTTPPrint(r *http.Request, resp *http.Response, target string) {
	fmt.Printf("[%s] [%s] [HTTP REQUEST] %s %s\nRequestSender: %s\nRequestTarget: %s\n", time.Now().Format("15:04:05"), l.Category, r.Method, r.URL, r.RemoteAddr, target)
	fmt.Printf("Headers:\n")
	for key, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	if r.Body != nil {
		body, _ := io.ReadAll(r.Body)
		fmt.Printf("Body: %s\n", string(body))
		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	if resp != nil {
		fmt.Printf("[%s] [%s] [HTTP RESPONSE] Status: %d %s\n", time.Now().Format("15:04:05"), l.Category, resp.StatusCode, resp.Status)
		fmt.Printf("Headers:\n")
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}

		if resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Body: %s\n", string(body))
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	} else {
		fmt.Printf("[%s] [%s] [HTTP RESPONSE] No response received\n", time.Now().Format("15:04:05"), l.Category)
	}

}
