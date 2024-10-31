package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type DefaultHttpG struct {
	Targets      []*TargetGroup
	ServerConfig *ServerConfig
}

func CreateDefaultHttpG(srv *ServerConfig, trg []*TargetGroup) *DefaultHttpG {
	if srv == nil || trg == nil {
		panic("nil config")
	}

	for _, tg := range trg {
		if tg.Targets == nil {
			panic("nil targets")
		}
		tg.Targets.Balancer = ChooseStrat(tg.BalanceStrategy)
	}
	return &DefaultHttpG{
		Targets:      trg,
		ServerConfig: srv,
	}
}

func (dg *DefaultHttpG) Serve() error {
	srv := http.Server{
		Addr:              dg.ServerConfig.ListenAddr,
		Handler:           dg.createRouters(),
		ReadTimeout:       time.Duration(dg.ServerConfig.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(dg.ServerConfig.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(dg.ServerConfig.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(dg.ServerConfig.IdleTimeout) * time.Second,
	}

	return srv.ListenAndServe()
}

func (dg *DefaultHttpG) createRouters() *http.ServeMux {
	mux := http.NewServeMux()

	for _, target := range dg.Targets {
		mux.HandleFunc(target.HandlePath, func(writer http.ResponseWriter, request *http.Request) {
			dg.hahandler(context.Background(), writer, request, target)
		})
		fmt.Println("created handle path: ", target.HandlePath, "for: ", target.Targets.URLs)
	}

	return mux
}

func (dg *DefaultHttpG) hahandler(ctx context.Context, w http.ResponseWriter, r *http.Request, targets *TargetGroup) {
	var client = http.Client{
		Timeout: 7 * time.Second,
	}

	respChan := make(chan *http.Response)
	errChan := make(chan error)

	c, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	go func() {
		defer close(respChan)
		defer close(errChan)

		t, err := targets.Targets.Balancer.PickServer(targets.Targets)
		if err != nil {
			errChan <- err
			return
		}

		req, er := http.NewRequestWithContext(c, r.Method, t.String(), r.Body)
		if er != nil {
			errChan <- err
			return
		}
		req.Header = r.Header

		resp, e := client.Do(req)
		if e != nil {
			errChan <- e
			return
		}

		respChan <- resp
	}()

	select {
	case <-ctx.Done():
		fmt.Println("context is done")
		http.Error(w, ctx.Err().Error(), http.StatusServiceUnavailable)
	case resp := <-respChan:
		resBytes, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("response: \n", string(resBytes))
		w.Write(resBytes)
	case e := <-errChan:
		log.Println(e)
		http.Error(w, e.Error(), http.StatusServiceUnavailable)
	}
}
