package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	robot "github.com/go-vgo/robotgo"
)

type Point struct {
	x int
	y int
}

var (
	counter = 1
	//点击链接的位置
	LinkPoint  = Point{1638, 608}
	ClosePoint = Point{1251, 163}
	InputPoint = Point{1372, 768}
	client     *http.Client

	stop bool
)

func initFlag() {
	flag.Parse()
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
	return _urls[rand.Intn(len(_urls))]
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

var (
	_urls = []string{
		"http://mp.weixin.qq.com/s?__biz=MzI2MzMxNzEzNA==&mid=2247484076&idx=1&sn=2b4b1dd2001d525e08966be9198d3f8d&scene=2&srcid=0804vRsESdgOtdaWTx4CET9Y&from=timeline&isappinstalled=0#wechat_redirect",

		"http://mp.weixin.qq.com/s?__biz=MjM5MDMyMzg2MA==&mid=2655481188&idx=1&sn=b22c1b7089ef132724d5c35b82372cb3&chksm=bdf5489f8a82c1895c72f488602a1f09b2ba7821a61d2499bf6ad35c88dcfa737af6675065a7#rd",

		"http://mp.weixin.qq.com/s?__biz=MjM5MDMyMzg2MA==&mid=2655481311&idx=4&sn=96a00896b5c800325ba5fc978ec15614&chksm=bdf548248a82c13223989f9469219a040d53c145eb9db4c309a6f4dfb1e4ed1a76dcc22b7196#rd",
	}
)
