package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"time"
)

const fn = "current_config.json"

type SSRGateServer struct {
	url string
	config *ssrConfig
	configChan chan *ssrConfig
	goodKeyWord string
	badKeyWord string
}

func newSSRGateServer(ssrUrl string, port int, goodKeyWord string, badKeyWord string) *SSRGateServer {
	serv := new(SSRGateServer)
	serv.url = ssrUrl
	serv.configChan = make(chan *ssrConfig)
	serv.goodKeyWord = goodKeyWord
	serv.badKeyWord = badKeyWord
	go runGost(serv.configChan, port)

	return serv
}

func (s *SSRGateServer) update() {
	cfgs, err := goodWayFromUrl(s.url, s.goodKeyWord, s.badKeyWord)
	if err != nil {
		logs(err)
		return
	}

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
		serv := newSSRGateServer(*u, *l, *k, *b)
		serv.Run()
	}
}
