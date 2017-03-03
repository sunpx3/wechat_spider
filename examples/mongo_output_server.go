package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	spider "github.com/sundy-li/wechat_spider"
	"gopkg.in/mgo.v2"
)

func main() {
	var port = "8899"
	spider.InitConfig(&spider.Config{
		// open it see detail logs
		Verbose: false,
		// Auto scroll
		AutoScroll: false,
	})
	spider.Regist(&CustomProcessor{})
	spider.Run(port)
}

//Just to implement Output Method of interface{} Processor
type CustomProcessor struct {
	spider.BaseProcessor
}

type M map[string]interface{}

func (c *CustomProcessor) Output() {
	for _, result := range c.UrlResults() {
		resp, err := http.Get(result.Url)
		if err != nil {
			println(err.Error())
			continue
		}
		uri, _ := url.ParseRequestURI(result.Url)
		id := md5Encode(uri.Query().Get("__biz") + "_" + uri.Query().Get("mid") + uri.Query().Get("idx"))
		bs, _ := ioutil.ReadAll(resp.Body)
		db.C("wx_article").UpsertId(id, M{"url": result.Url, "data": string(bs)})
	}
}

func init() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	db = session.DB("wechat") //数据库名称
}

func md5Encode(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

var (
	db *mgo.Database
)
