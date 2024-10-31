package test

import (
	"fmt"
	"loadBalancer/internal"
	"net/url"
	"testing"
)

func TestLeastConn(t *testing.T) {
	servers := []*internal.Server{
		{URL: &url.URL{Scheme: "http", Host: "localhost:5000", Path: "/api"}, ActiveConn: 10},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5001", Path: "/api"}, ActiveConn: 22},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5002", Path: "/api"}, ActiveConn: 5},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5003", Path: "/api"}, ActiveConn: 3},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5004", Path: "/api"}, ActiveConn: 4},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5005", Path: "/api"}, ActiveConn: 5},
	}

	strat := internal.ChooseStrat("leastConnections")

	fmt.Println("strat:", strat.ShowStrat())

	srv := strat.PickServer(servers)
	fmt.Println("choose: ", srv.URL)

	if srv.URL.String() != "http://localhost:5003/api" {
		t.Fatal(srv.URL.String())
	}
}

func TestRandom(t *testing.T) {
	servers := []*internal.Server{
		{URL: &url.URL{Scheme: "http", Host: "localhost:5000", Path: "/api"}},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5001", Path: "/api"}},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5002", Path: "/api"}},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5003", Path: "/api"}},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5004", Path: "/api"}},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5005", Path: "/api"}},
	}

	strat := internal.ChooseStrat("random")

	fmt.Println("strat: ", strat.ShowStrat())

	for i := 0; i < 12; i++ {
		srv := strat.PickServer(servers)
		fmt.Println("srv: ", srv.URL)
	}
}

func TestRoundRobin(t *testing.T) {
	servers := []*internal.Server{
		{URL: &url.URL{Scheme: "http", Host: "localhost:5000", Path: "/api"}},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5001", Path: "/api"}},
		{URL: &url.URL{Scheme: "http", Host: "localhost:5002", Path: "/api"}},
	}

	strat := internal.ChooseStrat("roundRobin")

	fmt.Println("strat: ", strat.ShowStrat())

	s1 := strat.PickServer(servers)
	s2 := strat.PickServer(servers)
	s3 := strat.PickServer(servers)
	s11 := strat.PickServer(servers)

	fmt.Println("s1 ", s1.URL)
	fmt.Println("s2 ", s2.URL)
	fmt.Println("s3 ", s3.URL)
	fmt.Println("s11 ", s11.URL)
	if s1.URL != servers[0].URL && s2.URL != servers[1].URL && s3.URL != servers[2].URL {
		t.Fatal("not match")
	}

	if s1.URL != s11.URL {
		t.Fatal("not match s1 s11")
	}
}
