package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const fn = "current_config.json"

var checkMethod = "tcp"
var checkAllMethod = "tcp"

var updateChan = make(chan bool)

type SSRGateServer struct {
	url          string
	configs      ssrConfigSlice
	configIndex  int
	config       *ssrConfig
	configChan   chan *ssrConfig
	goodKeyWords []string
	badKeyWords  []string
	port         int
}

func newSSRGateServer(ssrUrl string, port int, goodKeyWords []string, badKeyWords []string) *SSRGateServer {
	serv := new(SSRGateServer)
	serv.url = ssrUrl
	serv.configChan = make(chan *ssrConfig)
	serv.goodKeyWords = goodKeyWords
	serv.badKeyWords = badKeyWords
	serv.port = port
	go runSSR(serv.configChan, port)

	return serv
}

func (s *SSRGateServer) update() {
	if len(s.configs) > 0 && s.configIndex + 1 < len(s.configs) {
		s.configIndex += 1

	} else {
		s.configs = goodWays(s.url, s.goodKeyWords, s.badKeyWords)
		s.configIndex = 0
	}

	//s.config = bestWay(s.url, s.goodKeyWords, s.badKeyWords)
	s.config = s.configs[s.configIndex]
	s.configChan <- s.config

	j, _ := json.MarshalIndent(s.config, "", "  ")
	_ = ioutil.WriteFile(fn, j, 0644)
	time.Sleep(time.Second)
	s.check()
}

func (s *SSRGateServer) check() {
	log(s.config.Remarks, "check...")
	s.config.tcpPing()
	log(s.config.Remarks, s.config.Host, s.config.Port, "tcp ttl:", s.config.Ttl)
	if s.config.Ttl <= 0 {
		s.update()

	} else if checkMethod == "http" {
		ttl := HttpPing(s.port)
		log(s.config.Remarks, "http ttl:", ttl)
		if ttl < 0 {
			s.update()
		}
	}
}

func (s *SSRGateServer) updateBySignal(ch <-chan bool) {
	for {
		<- ch
		s.update()
	}
}

func (s *SSRGateServer) Run() {
	if pathExist(fn) {
		b, _ := ioutil.ReadFile(fn)
		cfg := new(ssrConfig)
		err := json.Unmarshal(b, cfg)
		check(err)

		s.config = cfg
		log("Read config from file:")
		logb(s.config)
		s.configChan <- s.config
		time.Sleep(time.Second)
		s.check()
	} else {
		s.update()
	}

	go s.updateBySignal(updateChan)
	for {
		time.Sleep(time.Second * 20)
		s.check()
	}
}

func httpHandle(w http.ResponseWriter, _ *http.Request) {
	updateChan <- true
	fmt.Fprintf(w, "SSR GATE UPDATE")
}

func main() {
	h := flag.Bool("h", false, "help")
	u := flag.String("u", "", "ssr url")
	l := flag.Int("l", 1080, "listen port")
	k := flag.String("k", "", "remarks match keywords")
	b := flag.String("b", "", "remarks match bad keywords")
	c := flag.String("c", "tcp", "the method that check node: [tcp|http]")
	m := flag.String("m", "tcp", "the method that check all nodes: [tcp|http]")

	flag.Parse()
	if *h || *u == "" {
		flag.Usage()
	} else {
		goodKeyWords := strings.Split(*k, " ")
		badKeyWords := strings.Split(*b, " ")

		log("good key words:", goodKeyWords)
		log("bad key words:", badKeyWords)
		serv := newSSRGateServer(*u, *l, goodKeyWords, badKeyWords)

		checkMethod = *c
		checkAllMethod = *m

		http.HandleFunc("/", httpHandle)
		go http.ListenAndServe(":8094", nil)

		serv.Run()
	}
}
