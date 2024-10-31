package main

import (
	"bytes"
	"encoding/json"
	"io"
	"loadBalancer/internal"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type Body struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Phone   string   `json:"phone"`
	Choose  bool     `json:"choose"`
	Params  string   `json:"params"`
	Prams   []string `json:"prams"`
}

func BenchmarkDefaultHttp_Serve(b *testing.B) {
	b.ReportAllocs()

	urlStr := "http://localhost:8080/api"

	u, _ := url.Parse(urlStr)

	dat := &Body{
		ID:      123142,
		Name:    "dsdckjcnck",
		Surname: "dsdck",
		Phone:   "+3130810942840944",
		Choose:  true,
		Params:  "ncndkvbciuwecuecubbclewbliu98y dxiuhdnddjxdmmd",
		Prams:   []string{"first", "second", "third", "fours"},
	}

	body, err := json.Marshal(dat)
	if err != nil {
		b.Error(err)
	}

	req := http.Request{
		Method: "POST",
		URL:    u,
		Header: http.Header{},
		Body:   io.NopCloser(bytes.NewReader(body)),
	}

	ss := internal.InitServerSide()

	if ss == nil {
		log.Fatal("init server failed")
	}

	infra := internal.InitInfra(ss)
	infra.PrintConf()

	go func() {
		err = infra.ServerSide.RunServer()
		if err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(1 * time.Second)

	for i := 0; i < 100; i++ {
		_, err := http.DefaultClient.Do(&req)
		if err != nil {
			b.Error(err)
		}
	}
}
