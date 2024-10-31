package internal

import (
	"flag"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type ServerSide struct {
	Choose string `yaml:"Choose"`

	FastHttp     *FastHttp     `yaml:"-"`
	DefaultHttp  *DefaultHttp  `yaml:"-"`
	DefaultHttpG *DefaultHttpG `yaml:"-"`
}

func (s *ServerSide) RunServer() error {
	switch s.Choose {
	case "fast":
		return s.FastHttp.Serve()
	case "default":
		return s.DefaultHttp.Serve()
	case "defaultG":
		return s.DefaultHttpG.Serve()
	default:
		panic("no server to run")
	}
	return nil
}

type ServerConfig struct {
	ListenAddr        string `yaml:"ListenAddr"`
	ReadTimeout       uint   `yaml:"ReadTimeout"`
	ReadHeaderTimeout uint   `yaml:"ReadHeaderTimeout"`
	WriteTimeout      uint   `yaml:"WriteTimeout"`
	IdleTimeout       uint   `yaml:"IdleTimeout"`
}

func InitServerSide() *ServerSide {
	choose := flag.String("choose", "fast", "Choose server")
	srvFile := flag.String("serverConf", "server.yaml", "Config file")
	trgFile := flag.String("trgFile", "trg.yaml", "Config file")
	flag.Parse()

	srv, trg, err := readServerConfig(*srvFile, *trgFile)
	if err != nil {
		log.Fatal(err)
	}

	switch *choose {
	case "fast":
		return &ServerSide{
			Choose:       *choose,
			FastHttp:     CreateFastHttp(srv, trg),
			DefaultHttp:  nil,
			DefaultHttpG: nil,
		}
	case "default":
		return &ServerSide{
			Choose:       *choose,
			FastHttp:     nil,
			DefaultHttp:  CreateDefaultHttp(srv, trg),
			DefaultHttpG: nil,
		}
	case "defaultG":
		return &ServerSide{
			Choose:       *choose,
			FastHttp:     nil,
			DefaultHttp:  nil,
			DefaultHttpG: CreateDefaultHttpG(srv, trg),
		}
	default:
		panic("no server selected")
	}
	return &ServerSide{}
}

type STGWrapper struct {
	ServConf *ServerConfig
	TrgConf  []*TargetGroup
}

func readServerConfig(srvPath, trgPath string) (*ServerConfig, []*TargetGroup, error) {
	var wrp STGWrapper
	wrp.ServConf = &ServerConfig{}
	wrp.TrgConf = []*TargetGroup{}

	srvFile, err := os.OpenFile(srvPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}
	defer srvFile.Close()

	trgFile, e := os.OpenFile(trgPath, os.O_RDONLY, os.ModePerm)
	if e != nil {
		return nil, nil, e
	}
	defer trgFile.Close()

	err = yaml.NewDecoder(trgFile).Decode(&wrp.TrgConf)
	if err != nil {
		return nil, nil, err
	}

	err = yaml.NewDecoder(srvFile).Decode(wrp.ServConf)
	if err != nil {
		return nil, nil, err
	}

	return wrp.ServConf, wrp.TrgConf, nil
}
