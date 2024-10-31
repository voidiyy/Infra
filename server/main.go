package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	port := flag.String("port", "8080", "port lala")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/db", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		var buff bytes.Buffer
		_, err := buff.ReadFrom(request.Body)
		if err != nil {
			http.Error(writer, "Failed to read request body", http.StatusBadRequest)
			return
		}

		fmt.Println("accept request: ", request.RequestURI, "method: ", request.Method)
		fmt.Println("body length: ", len(buff.Bytes()))
		fmt.Println("___________________")
		writer.Write([]byte("DB response"))
	})

	mux.HandleFunc("/grpc", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		var buff bytes.Buffer
		_, err := buff.ReadFrom(request.Body)
		if err != nil {
			http.Error(writer, "Failed to read request body", http.StatusBadRequest)
			return
		}
		fmt.Println("accept request: ", request.RequestURI, "method: ", request.Method)
		fmt.Println("body length: ", len(buff.Bytes()))
		fmt.Println("___________________")
		writer.Write([]byte("GRPC response"))
	})

	mux.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		var buff bytes.Buffer
		_, err := buff.ReadFrom(request.Body)
		if err != nil {
			http.Error(writer, "Failed to read request body", http.StatusBadRequest)
			return
		}

		fmt.Println("accept request: ", request.RequestURI, "method: ", request.Method)
		fmt.Println("body length: ", len(buff.Bytes()))
		fmt.Println("___________________")

		_, _ = writer.Write([]byte("API response"))
	})

	fmt.Println("listen: ", *port)
	err := http.ListenAndServe(":"+*port, mux)
	if err != nil {
		fmt.Println(err)
	}
}

func printRequest(req *http.Request, body *io.ReadCloser) {

	var buf bytes.Buffer
	_, err := buf.ReadFrom(*body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[+]")
	fmt.Println("REQUEST")
	fmt.Println("From: ", req.RemoteAddr)
	fmt.Println("Headers: ")
	for k, v := range req.Header {
		fmt.Println("\t", k, ":", v)
	}
	fmt.Println("Body: \n", "\t", buf.String())

	reqBytes, _ := httputil.DumpRequestOut(req, true)

	fmt.Println(string(reqBytes))
	fmt.Println("__________________________________")

	return
}
