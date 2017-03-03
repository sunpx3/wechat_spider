package main

import (
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	spider "github.com/sundy-li/wechat_spider"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

/**
	create database wechat;
	create table wx_article(id bigint(15) not null primary key , url varchar(512), data text, update_time bigint(15));
**/

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
		id := hash(uri.Query().Get("__biz") + "_" + uri.Query().Get("mid") + uri.Query().Get("idx"))
		bs, _ := ioutil.ReadAll(resp.Body)

		stmt, err := db.Prepare("insert ignore into wx_article(id,url,data,update_time) values(?,?,?,?)")
		if err != nil {
			println(err.Error())
			continue
		}
		_, err = stmt.Exec(id, result.Url, string(bs), time.Now().Unix())
		if err != nil {
			println(err.Error())
		}
	}
}

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:4000)/wechat?charset=utf8")
	if err != nil {
		panic(err)
	}
}

func hash(s string) int64 {
	h := fnv.New32()
	h.Write([]byte(s))
	return int64(h.Sum32())
}

var (
	db *sql.DB
)
