package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"strings"
	"time"
)

const fn = "current_config.json"

type SSRGateServer struct {
	url string
	config *ssrConfig
	configChan chan *ssrConfig
	goodKeyWords []string
	badKeyWords []string
}

func newSSRGateServer(ssrUrl string, port int, goodKeyWords []string, badKeyWords []string) *SSRGateServer {
	serv := new(SSRGateServer)
	serv.url = ssrUrl
	serv.configChan = make(chan *ssrConfig)
	serv.goodKeyWords = goodKeyWords
	serv.badKeyWords = badKeyWords
	go runSSR(serv.configChan, port)

	return serv
}

func (s *SSRGateServer) update() {
	cfgs := goodWaysFromUrl(s.url, s.goodKeyWords, s.badKeyWords)

	s.config = bestWay(cfgs)
	s.configChan <- s.config

	j, _ := json.MarshalIndent(s.config, "", "  ")
	_ = ioutil.WriteFile(fn, j, 0644)
}

func (s *SSRGateServer) check() {
	logs("check...")
	s.config.ping()
	logs(s.config.Host, s.config.Port, s.config.Ttl)
	if s.config.Ttl <= 0 {
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
		logs("Read config from file:")
		logb(s.config)
		s.configChan <- s.config
		s.check()
	} else {
		s.update()
	}

	for {
		time.Sleep(time.Second * 30)
		s.check()
	}
}

func main() {
	h := flag.Bool("h", false, "help")
	u := flag.String("u", "", "ssr url")
	l := flag.Int("l", 1080, "listen port")
	k := flag.String("k", "", "remarks match keywords")
	b := flag.String("b", "", "remarks match bad keywords")

	flag.Parse()
	if *h || *u == "" {
		flag.Usage()
	} else {
		goodKeyWords := strings.Split(*k, " ")
		badKeyWords := strings.Split(*b, " ")
		serv := newSSRGateServer(*u, *l, goodKeyWords, badKeyWords)
		serv.Run()
	}
}
