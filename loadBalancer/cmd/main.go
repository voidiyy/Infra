package main

import (
	"loadBalancer/internal"
	"log"
	_ "net/http/pprof"
)

func main() {

	ss := internal.InitServerSide()

	if ss == nil {
		log.Fatal("init server failed")
	}

	infra := internal.InitInfra(ss)
	infra.PrintConf()
	err := infra.ServerSide.RunServer()
	if err != nil {
		panic(err)
	}
}
