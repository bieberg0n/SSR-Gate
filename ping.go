package main

import (
	"golang.org/x/net/proxy"
	"net"
	"net/http"
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

func HttpPing(socksPort int) bool {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+strconv.Itoa(socksPort), nil, &net.Dialer{Timeout: 3 * time.Second})
	//KeepAlive: 10 * time.Second,
	if err != nil {
		log("get dialer error", dialer)
		return false
	}
	httpTransport := &http.Transport{Dial: dialer.Dial}
	httpClient := &http.Client{Transport: httpTransport}
	resp, err := httpClient.Get("https://www.google.com/")
	if err != nil {
		log("http get error:", err)
		return false
	} else {
		defer resp.Body.Close()
		return true
	}
}