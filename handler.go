package wechat_spider

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/elazarl/goproxy"
)

var (
	Logger = log.New(os.Stderr, "", log.LstdFlags)
)

func ProxyHandle(proc Processor) func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp == nil || resp.StatusCode != 200 {
			return resp
		}

		if ctx.Req.URL.Path == `/mp/getmasssendmsg` || (ctx.Req.URL.Path == `/mp/profile_ext` && ctx.Req.URL.Query().Get("action") == "home") {
			handleList(resp, ctx, proc)
			header = ctx.Req.Header
		} else if ctx.Req.URL.Path == `/s` {
			header = ctx.Req.Header
			handleDetail(resp, ctx, proc)
		} else if ctx.Req.URL.Path == "/mp/getappmsgext" {
			handleMetrics(resp, ctx, proc)
		} else {
			return resp
		}
		if rootConfig.Verbose {
			Logger.Printf("Hijacked of %s:%s", ctx.Req.Method, ctx.Req.URL.RequestURI())
		}
		return resp
	}
}

func handleList(resp *http.Response, ctx *goproxy.ProxyCtx, proc Processor) {
	var err error
	p := getProcessor(ctx.Req, proc)
	data, err := p.ProcessList(resp, ctx)
	if err != nil {
		log.Println(err.Error())
	}
	go func() {
		if rootConfig.AutoScroll {
			p.ProcessPages()
		}
		p.Output()
	}()
	var buf = bytes.NewBuffer(data)
	resp.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
}

func handleDetail(resp *http.Response, ctx *goproxy.ProxyCtx, proc Processor) {
	var err error
	p := getProcessor(ctx.Req, proc)
	data, err := p.ProcessDetail(resp, ctx)
	if err != nil {
		log.Println(err.Error())
	}
	var buf = bytes.NewBuffer(data)
	go p.Output()
	resp.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
}

func handleMetrics(resp *http.Response, ctx *goproxy.ProxyCtx, proc Processor) {
	var err error
	p := getProcessor(ctx.Req, proc)
	data, err := p.ProcessMetrics(resp, ctx)
	if err != nil {
		Logger.Println(err.Error())
	}
	go p.Output()
	resp.Body = ioutil.NopCloser(bytes.NewReader(data))
}

// get the processor from cache
func getProcessor(req *http.Request, proc Processor) Processor {
	t := reflect.TypeOf(proc)
	v := reflect.New(t.Elem())
	p := v.Interface().(Processor)
	return p
}
