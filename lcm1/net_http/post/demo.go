package main

import (
	"flag"
	"io/ioutil"
	"lcm1/net_http/post/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

// InputTestNum 瞬间访问峰值测试
var InputTestNum = flag.Int("num", 1, "input test num")

// InputPos input body pos pos 对应测试接口，由于简单写不做容错处理，也就需要自己判断该接口（config.go）有没有添加
var InputPos = flag.Int("pos", 7, "input body pos")

func httpPost(id int, pos config.Bodypos) {
	resp, err := http.Post(config.Urls[pos],
		"raw",
		strings.NewReader(config.Bodys[pos]))

	if err != nil {
		log.Printf("[%d] 1 %v\n", id, err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("[%d] 2 %v\n", id, err)
		return
	}

	log.Printf("[%d] body:%s\n", id, string(body))
}

func main() {
	flag.Parse()
	pos := config.Bodypos(*InputPos)
	testLen := *InputTestNum
	for i := 0; i != testLen; i++ {
		go httpPost(i, pos)
	}

	test := make(chan os.Signal, 1)
	signal.Notify(test, os.Interrupt, os.Kill)
	sig := <-test
	log.Printf("webTest closing down (signal: %v)", sig)
}
