package wechat_spider

import (
	"log"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

type spider struct {
	proxy *goproxy.ProxyHttpServer
}

var _spider = newSpider()

func Regist(proc Processor) {
	_spider.Regist(proc)
}

func Run(port string) {
	_spider.Run(port)
}

func newSpider() *spider {
	sp := &spider{}
	sp.proxy = goproxy.NewProxyHttpServer()

	sp.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	return sp
}

func (s *spider) Regist(proc Processor) {
	s.proxy.OnResponse().DoFunc(ProxyHandle(proc))
}

func (s *spider) Run(port string) {
	if rootConfig.Compress {
		s.proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			host := req.URL.Host
			if !strings.Contains(host, "mp.weixin.qq.com") {
				req, _ = http.NewRequest("GET", "http://mp.weixin.qq.com/notfound", nil)
				return req, nil
			}
			return req, nil
		})
	}
	log.Println("server will at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, s.proxy))
}

var (
	rootConfig = &Config{
		Verbose:    false,
		AutoScroll: false,
		Compress:   true,
	}
)

type Config struct {
	Verbose    bool // Debug
	AutoScroll bool // Auto scroll page to hijack all history articles
	Compress   bool // To ingore other request just to save the net cost
}

func InitConfig(conf *Config) {
	rootConfig = conf
}

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}
