package main

import (
	"fmt"

	spider "github.com/sundy-li/wechat_spider"
)

func main() {
	var port = "8899"
	spider.InitConfig(&spider.Config{
		Verbose:    false, // Open to see detail logs
		AutoScroll: false, // Open to crawl scroll pages
	})
	spider.Regist(&CustomProcessor{})
	spider.Run(port)
}

//Just to implement Output Method of interface{} Processor
type CustomProcessor struct {
	spider.BaseProcessor
}

func (c *CustomProcessor) Output() {
	switch c.Type {
	case spider.TypeList:
		//do nothing
		fmt.Printf("url size ==> %#v\n", len(c.UrlResults()))
	case spider.TypeDetail:
		fmt.Printf("url %s %s is being spidered\n", c.DetailResult().Id, c.DetailResult().Url)
	case spider.TypeMetric:
		fmt.Printf("url %s %s metric %#v is being spidered\n", c.DetailResult().Id, c.DetailResult().Url, c.DetailResult().Appmsgstat)
	}
}
