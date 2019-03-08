package main

import (
	"time"
	"net"
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
			logs(err)
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
