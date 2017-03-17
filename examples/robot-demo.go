package main

import (
	robot "github.com/go-vgo/robotgo"
	"fmt"
	"time"
	"log"
	"net/http"
	"io/ioutil"
	"net/url"
	"flag"
	"encoding/json"
)

type Point struct {
	x int
	y int
}

var (
	counter = 1
	//点击链接的位置
	LinkPoint = Point{1638, 608}
	
	ClosePoint = Point{1251, 163}
	InputPoint = Point{1372, 768}
	client     *http.Client
	
	stop bool = false
	
	proxy string = "http://10.0.0.134:8899"
	
	nowUrl string
)

func initFlag() {
	flag.StringVar(&proxy, "proxy", "http://10.0.0.134:8899", "run port")
	flag.Parse()
	fmt.Println("代理地址：", proxy)
}

func initClient() {
	u, err := url.Parse(proxy)
	if err != nil {
		//log.Println("err to parse proxy url")
		panic(err)
		//panic(fmt.Sprintf("err to parse proxy url: %s", proxy))
	}
	f := http.ProxyURL(u)
	transport := &http.Transport{Proxy: f}
	client = &http.Client{Transport: transport}
}

func initPoint() {
	fmt.Println("输入点击链接的位置")
	sleep(2)
	x, y := robot.GetMousePos()
	LinkPoint = Point{x, y}
	fmt.Println(LinkPoint)
	
	//fmt.Println("输入点击关闭微信文章窗口的位置")
	//sleep(2)
	//x, y = robot.GetMousePos()
	//ClosePoint = Point{x, y}
	//fmt.Println(ClosePoint)
	
	fmt.Println("输入输入框的位置")
	sleep(2)
	x, y = robot.GetMousePos()
	InputPoint = Point{x, y}
	fmt.Println(InputPoint)
	fmt.Println("点击输入框")
	fmt.Println("点击F8退出")
	sleep(2)
}

// n second
func sleep(n int) {
	time.Sleep(time.Duration(n) * time.Second)
}

func NextUrl() string {
	return getNextUrl()
}

type Article struct {
	Url       string `json:"url"`
	Id        string `json:"_id"`
	Title     string `json:"title"`
	UpdateTme int64 `json:"update_tme"`
}

func getNextUrl() string {
	req, err := http.NewRequest("GET", "http://mp.weixin.qq.com/proxy/nexturl?url="+nowUrl, nil)
	if err != nil {
		log.Println("init next url req err")
		return ""
	}
	
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		log.Println("get next request err")
		return ""
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read resp err")
		return ""
	}
	
	article := &Article{}
	json.Unmarshal(body, article)
	
	return article.Url
}

func stopRobot() {
	key := robot.AddEvent("f8")
	if key == 0 {
		log.Println("stop")
		stop = true
	}
}

func process() {
	fmt.Printf("requested: %d %s \n", counter, time.Now().String())
	
	u := NextUrl()
	if u == "" {
		sleep(10)
		return
	}
	robot.TypeString(u)
	
	robot.KeyTap("enter")
	
	robot.MoveMouseSmooth(LinkPoint.x, LinkPoint.y, 0.5)
	robot.Click(LinkPoint.x, LinkPoint.y)
	//
	//robot.MoveMouseSmooth(ClosePoint.x, ClosePoint.y, 0.5)
	//robot.Click(ClosePoint.x, ClosePoint.y)
	//
	robot.MoveMouseSmooth(InputPoint.x, InputPoint.y, 0.5)
	robot.Click(InputPoint.x, InputPoint.y)
	
	counter = counter + 1
}


func init() {
	initFlag()
	initClient()
	initPoint()
}

func start() {
	go stopRobot()
	for !stop {
		process()
		time.Sleep(time.Second * 5)
	}
}

func main() {
	start()
}
