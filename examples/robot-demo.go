package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
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

	stop        bool
	sleepSecond int
)

func initFlag() {
	flag.IntVar(&sleepSecond, "s", 5, "sleep second")
	flag.Parse()
}

func initPoint() {
	fmt.Println("鼠标移动到点击链接的位置(f3结束)")
	key := robot.AddEvent("f3")
	sleep(1)
	if key == 0 {
		x, y := robot.GetMousePos()
		LinkPoint = Point{x, y}
		fmt.Println("链接地址", x, y)
	}
	fmt.Println("鼠标移动到输入框的位置(f3结束)")
	key = robot.AddEvent("f3")
	sleep(1)
	if key == 0 {
		x, y := robot.GetMousePos()
		InputPoint = Point{x, y}
		fmt.Println("链接地址", x, y)
	}
	fmt.Println("程序开始,点击F8退出")
}

func stopRobot() {
	for {
		key := robot.AddEvent("f8")
		sleep(1)
		if key == 0 {
			stop = !stop
			if stop {
				log.Println("stopped")
			} else {
				log.Println("started")
			}
		}
	}
}

func process() {
	fmt.Printf("requested: %d %s \n", counter, time.Now().String())

	//这里实现批量化即可
	nextUrl := NextUrl()
	if nextUrl == "" {
		return
	}
	robot.TypeString(nextUrl)
	robot.KeyTap("enter")
	robot.MoveMouseSmooth(LinkPoint.x, LinkPoint.y, 0.5)
	robot.Click(LinkPoint.x, LinkPoint.y)

	robot.MoveMouseSmooth(InputPoint.x, InputPoint.y, 0.5)
	robot.Click(InputPoint.x, InputPoint.y)

	counter = counter + 1
}

func init() {
	initFlag()
	initPoint()
}

func main() {
	go stopRobot()
	for {
		if stop {
			time.Sleep(1 * time.Second)
			continue
		}
		process()
		sleep(sleepSecond)
	}
}

// n second
func sleep(n int) {
	time.Sleep(time.Duration(n) * time.Second)
}

func NextUrl() string {
	return _urls[rand.Intn(len(_urls))]
}

var (
	_urls = []string{
		"http://mp.weixin.qq.com/s?__biz=MzI2MzMxNzEzNA==&mid=2247484076&idx=1&sn=2b4b1dd2001d525e08966be9198d3f8d&scene=2&srcid=0804vRsESdgOtdaWTx4CET9Y&from=timeline&isappinstalled=0#wechat_redirect",

		"http://mp.weixin.qq.com/s?__biz=MjM5MDMyMzg2MA==&mid=2655481188&idx=1&sn=b22c1b7089ef132724d5c35b82372cb3&chksm=bdf5489f8a82c1895c72f488602a1f09b2ba7821a61d2499bf6ad35c88dcfa737af6675065a7#rd",

		"http://mp.weixin.qq.com/s?__biz=MjM5MDMyMzg2MA==&mid=2655481311&idx=4&sn=96a00896b5c800325ba5fc978ec15614&chksm=bdf548248a82c13223989f9469219a040d53c145eb9db4c309a6f4dfb1e4ed1a76dcc22b7196#rd",
	}
)
