package main

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ssrConfig struct {
	Host       string
	Port       int
	Protocol   string
	Method     string
	Obfs       string
	Password   string
	ObfsParam  string
	ProtoParam string
	Remarks    string
	Group      string
	Udpport    int
	Uot        bool
	Ttl        int
}

func (c *ssrConfig) ping () {
	c.Ttl = TcpPing(c.Host + ":" + strconv.Itoa(c.Port))
}

func parseSSRUrl (url string) (*ssrConfig, error) {
	url, err := b64decode(url[6:])
	if err != nil {
		return nil, err
	}

	ssrConfig := new(ssrConfig)

	parts := strings.Split(url, "/?")
	leftParts := strings.Split(parts[0], ":")
	rightParts := strings.Split(parts[1], "&")

	ssrConfig.Host = leftParts[0]
	ssrConfig.Port, err = strconv.Atoi(leftParts[1])
	if err != nil {
		return nil, err
	}

	ssrConfig.Protocol = leftParts[2]
	ssrConfig.Method = leftParts[3]
	ssrConfig.Obfs = leftParts[4]
	ssrConfig.Password, err = b64decode(leftParts[5])
	if err != nil {
		return nil, err
	}

	for _, param := range rightParts {
		kv := strings.Split(param, "=")
		switch kv[0] {
		case "obfsparam":
			ssrConfig.ObfsParam, err = b64decode(kv[1])
			if err != nil {
				return nil, err
			}
		case "protoparam":
			ssrConfig.ProtoParam, err = b64decode(kv[1])
			if err != nil {
				return nil, err
			}
		case "remarks":
			ssrConfig.Remarks, err = b64decode(kv[1])
			if err != nil {
				return nil, err
			}
		case "group":
			ssrConfig.Group, err = b64decode(kv[1])
			if err != nil {
				return nil, err
			}
		case "udpport":
			ssrConfig.Udpport, err = strconv.Atoi(kv[1])
			if err != nil {
				return nil, err
			}
		case "uot":
			if kv[1] == "0" {
				ssrConfig.Uot = false
			} else {
				ssrConfig.Uot = true
			}
		}
	}

	return ssrConfig, nil
}

func readSSR (data string) (map[string]*ssrConfig, error) {
	ssrs := strFilter(strings.Split(data, "\n"), func(s string) bool {
		return s != ""
	})

	cfgs := map[string]*ssrConfig{}
	for _, ssr := range ssrs {
		cfg, err := parseSSRUrl(ssr)
		if err != nil {
			logs(err, ":", ssr)
			continue
		}
		cfgs[cfg.Host] = cfg
	}
	return cfgs, nil
}

func cfgsFromUrl (url string) (map[string]*ssrConfig, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	decode, err := b64decode(buf.String())
	if err != nil {
		return nil, err
	}

	return readSSR(decode)
}

func bestWay(cfgs []*ssrConfig) (*ssrConfig) {
	ttlHostMap := map[int]*ssrConfig{}
	var ttls []int
	for _, cfg := range cfgs {
		ttlHostMap[cfg.Ttl] = cfg
		ttls = append(ttls, cfg.Ttl)
	}
	minTTL := minInt(ttls)
	best := ttlHostMap[minTTL]
	logs("best addr:", best.Remarks, best.Host, best.Port, minTTL)

	return ttlHostMap[minTTL]
}

func goodWays(cfgs map[string]*ssrConfig, goodKeyWords []string, badKeyWords []string) ([]*ssrConfig) {
	var goodCfgs []*ssrConfig
	for _, cfg := range cfgs {
		if (len(badKeyWords) != 0 && anyStrsInStr(cfg.Remarks, badKeyWords)) ||
			(len(goodKeyWords) != 0 && !anyStrsInStr(cfg.Remarks, goodKeyWords)) {
			logs(cfg.Host, cfg.Remarks, "BAN")
			continue
		}

		cfg.ping()
		logs(cfg.Host, cfg.Remarks, "ttl:", cfg.Ttl)
		if cfg.Ttl > 0 {
			goodCfgs = append(goodCfgs, cfg)
		}
	}
	return goodCfgs
}

func goodWayFromUrl (url string, goodKeyWords []string, badKeyWords []string) []*ssrConfig {
	info("http get ssr config...")
	for {
		cfgs, err := cfgsFromUrl(url)
		if err != nil {
			info("http get ssr config error:", err)
			time.Sleep(5 * time.Second)
		} else {
			return goodWays(cfgs, goodKeyWords, badKeyWords)
		}
	}

}
