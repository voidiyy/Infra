package internal

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type DefaultHttp struct {
	Targets      []*TargetGroup `yaml:"targets"`
	ServerConfig *ServerConfig  `yaml:"ServerConfig"`
}

func CreateDefaultHttp(srv *ServerConfig, trg []*TargetGroup) *DefaultHttp {
	var targets []*TargetGroup

	if srv == nil || trg == nil {
		panic("nil config")
	}

	for _, tg := range trg {
		targets = append(targets, tg)
	}

	return &DefaultHttp{
		Targets:      targets,
		ServerConfig: srv,
	}
}

func (d *DefaultHttp) Serve() error {
	srv := http.Server{
		Addr:              d.ServerConfig.ListenAddr,
		Handler:           d.createDefaultRouter(),
		ReadTimeout:       time.Duration(d.ServerConfig.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(d.ServerConfig.ReadHeaderTimeout) * time.Second,
		IdleTimeout:       time.Duration(d.ServerConfig.IdleTimeout) * time.Second,
		WriteTimeout:      time.Duration(d.ServerConfig.WriteTimeout) * time.Second,
	}

	fmt.Println("server run ... ", d.ServerConfig.ListenAddr)
	return srv.ListenAndServe()
}

func (d *DefaultHttp) createDefaultRouter() *http.ServeMux {
	mux := http.NewServeMux()

	for _, target := range d.Targets {
		mux.HandleFunc(target.HandlePath, func(writer http.ResponseWriter, request *http.Request) {
			d.Forward(writer, request, target)
		})
		fmt.Println("created router: ", target.HandlePath, "with: ", target.Targets.URLs)
	}

	return mux
}

func (d *DefaultHttp) Forward(writer http.ResponseWriter, request *http.Request, targets *TargetGroup) {

	var client = http.Client{
		Timeout: 7 * time.Second,
	}

	t, err := targets.Targets.Balancer.PickServer(targets.Targets)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	req, er := http.NewRequest(request.Method, t.String(), request.Body)
	if er != nil {
		http.Error(writer, er.Error(), http.StatusInternalServerError)
		return
	}
	req.Header = request.Header

	resp, e := client.Do(req)
	if e != nil {
		http.Error(writer, e.Error(), http.StatusInternalServerError)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	_, err = writer.Write(body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("--------------------")
	fmt.Println("request: ", req)
	fmt.Println("--------------------")
	fmt.Println("response: ", resp)
	fmt.Println("--------------------")
}
