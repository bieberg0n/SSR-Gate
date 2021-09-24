package main

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func tcpPingOnce(addr string) (int, error) {
	oldTime := time.Now()
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	ttl := time.Now().Sub(oldTime).Seconds()
	return int(ttl * 1000), nil
}

func TcpPing(addr string) int {
	var ttls []int
	for i := 1; i <= 3; i++ {
		ttl, err := tcpPingOnce(addr)
		if err != nil {
			log(err)
		} else {
			ttls = append(ttls, ttl)
		}
	}

	if len(ttls) == 0 {
		return -1
	} else {
		return minInt(ttls)
	}
}

func httpPingOnce(socksPort int) (int, error) {
	oldTime := time.Now()
	socksProxy := "socks5://127.0.0.1:" + strconv.Itoa(socksPort)
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(socksProxy)
	}

	httpTransport := &http.Transport{Proxy: proxy}
	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   3 * time.Second,
	}
	resp, err := httpClient.Head("https://www.google.com/")
	ttl := time.Now().Sub(oldTime).Seconds()
	if err != nil {
		log("http get error:", err)
		return 0, err
	}

	_ = resp.Body.Close()
	return int(ttl * 1000), nil
}

func HttpPing(socksPort int) int {
	var ttls []int
	for i := 0; i < 3; i++ {
		ttl, err := httpPingOnce(socksPort)
		if err != nil {
			log(err)
		} else {
			ttls = append(ttls, ttl)
		}
	}

	if len(ttls) == 0 {
		return -1
	} else {
		return minInt(ttls)
	}
}

func runAndHttpPing (cfg *ssrConfig) int {
	ch := make(chan *ssrConfig)
	port := int(RandInt64(30000, 60000))
	log("check:", cfg.Remarks, port)
	go runSSR(ch, "127.0.0.1", port)
	ch <- cfg
	time.Sleep(500 * time.Millisecond)
	ttl := HttpPing(port)
	log(cfg.Remarks, "ttl:", ttl)
	ch <- nil
	return ttl
}