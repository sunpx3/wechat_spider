package main

import (
	spider "github.com/sundy-li/wechat_spider"
)

func main() {
	var port = "8899"
	spider.InitConfig(&spider.Config{
		Verbose:    true,  // Open to see detail logs
		AutoScroll: true,  // Open to crawl scroll pages
		Compress:   false, // Ingore other request to save the
	})
	spider.Regist(spider.NewBaseProcessor())
	spider.Run(port)
}
