package internal

import (
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

type FastHttp struct {
	Targets      []*TargetGroup `yaml:"targets"`
	ServerConfig *ServerConfig  `yaml:"ServerConfig"`
}

func CreateFastHttp(srv *ServerConfig, trg []*TargetGroup) *FastHttp {
	if srv == nil || trg == nil {
		panic("nil config")
	}
	return &FastHttp{
		ServerConfig: srv,
		Targets:      trg,
	}
}

func (f *FastHttp) Serve() error {
	r := f.createFastRouter()

	srv := fasthttp.Server{
		ReadTimeout:  time.Duration(f.ServerConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(f.ServerConfig.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(f.ServerConfig.IdleTimeout) * time.Second,
	}
	fmt.Println("server started: ", f.ServerConfig.ListenAddr)

	srv.Handler = r.HandleRequest
	return srv.ListenAndServe(f.ServerConfig.ListenAddr)
}

func (f *FastHttp) createFastRouter() *routing.Router {
	router := routing.New()

	for _, target := range f.Targets {
		router.To("GET,POST", target.HandlePath, func(context *routing.Context) error {
			f.requestHandler(context.RequestCtx, target)
			return nil
		})
		fmt.Println("created router: ", target.HandlePath, "with: ", target.Targets.URLs)
	}

	return router
}

func (f *FastHttp) requestHandler(ctx *fasthttp.RequestCtx, target *TargetGroup) {
	client := fasthttp.Client{
		ReadTimeout:  time.Duration(f.ServerConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(f.ServerConfig.WriteTimeout) * time.Second,
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	trg, err := ChooseStrat(target.BalanceStrategy).PickServer(target.Targets)
	if err != nil || trg == nil {
		ctx.Error("invalid target", fasthttp.StatusServiceUnavailable)
		return
	}

	req.SetBody(ctx.PostBody())
	req.SetRequestURI(trg.RequestURI())
	req.SetHost(trg.Host)

	req.Header.SetMethodBytes(ctx.Method())
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = client.DoTimeout(req, resp, 8*time.Second)
	if err != nil {
		log.Printf("Error sending request to backend: %v", err)
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		return
	}

	resp.Header.VisitAll(func(key, value []byte) {
		ctx.Response.Header.Set(string(key), string(value))
	})
	ctx.Response.SetStatusCode(resp.StatusCode())
	ctx.Response.SetBody(resp.Body())

	fmt.Println("--------------------")
	fmt.Println("request: ", req)
	fmt.Println("--------------------")
	fmt.Println("response: ", resp)
	fmt.Println("--------------------")
}
