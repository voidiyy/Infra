package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
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
	HTTPPrint(r *http.Request, resp *http.Response, target string, args ...interface{})
}

func InitLogan(category string) *Log {
	switch category {
	case "app":
		return &Log{Category: category}
	case "server":
		return &Log{Category: category}
	case "http":
		return &Log{Category: category}
	default:
		return &Log{Category: "app"}
	}
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

func (l *Log) HTTPPrint(r *http.Request, resp *http.Response, target string, args ...interface{}) {
	// Читаємо та зберігаємо копію тіла запиту
	var requestBodyCopy bytes.Buffer
	if r.Body != nil {
		// Зберігаємо оригінал тіла
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading request body:", err)
		} else {
			// Створюємо нове тіло для подальшого використання
			requestBodyCopy.Write(bodyBytes)
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Відновлюємо тіло запиту для подальшого використання
		}
	}

	// Читаємо та зберігаємо копію тіла відповіді
	var responseBodyCopy bytes.Buffer
	if resp.Body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		} else {
			responseBodyCopy.Write(bodyBytes)
			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Відновлюємо тіло відповіді для подальшого використання
		}
	}

	// Логуємо запит
	requestDump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
	}

	// Логуємо відповідь
	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println("Error dumping response:", err)
	}

	// Виведення логів
	fmt.Println("Timestamp: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("Target: ", target)
	fmt.Println("MetaData: ", args)
	fmt.Println("Request: \n", string(requestDump))
	fmt.Println("Request Body: \n", requestBodyCopy.String())
	fmt.Println("Response: \n", string(responseDump))
	fmt.Println("Response Body: \n", responseBodyCopy.String())
}
