package internal

import (
	"fmt"
)

type Infra struct {
	ServerSide *ServerSide `yaml:"ServerSide"`

	lo Logan `yaml:"-"`
}

func InitInfra(ss *ServerSide) *Infra {
	return &Infra{
		ServerSide: ss,
		lo:         InitLogan(""),
	}
}

func (i *Infra) PrintConf() {
	fmt.Println("INFRA: ")
	fmt.Println("choose: ", i.ServerSide.Choose)

	switch i.ServerSide.Choose {
	case "fast":
		for _, target := range i.ServerSide.FastHttp.Targets {
			fmt.Println("___")
			fmt.Println("\tName: ", target.Name)
			fmt.Println("\tHandle path: ", target.HandlePath)
			fmt.Println("\tBalanceStrategy: ", target.BalanceStrategy)
			fmt.Println("\tURLs: ", target.Targets.URLs)
		}
	case "default":
		for _, target := range i.ServerSide.DefaultHttp.Targets {
			fmt.Println("___")
			fmt.Println("\tName: ", target.Name)
			fmt.Println("\tHandle path: ", target.HandlePath)
			fmt.Println("\tBalanceStrategy: ", target.BalanceStrategy)
			fmt.Println("\tURLs: ", target.Targets.URLs)
		}
	case "defaultG":
		for _, target := range i.ServerSide.DefaultHttpG.Targets {
			fmt.Println("___")
			fmt.Println("\tName: ", target.Name)
			fmt.Println("\tHandle path: ", target.HandlePath)
			fmt.Println("\tBalanceStrategy: ", target.BalanceStrategy)
			fmt.Println("\tURLs: ", target.Targets.URLs)
		}
	default:
		fmt.Println("Unknown choose: ", i.ServerSide.Choose)
	}
}
