package main

import (
	"bytes"
	"net/http"
	"sort"
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

func (c *ssrConfig) httpPing() {
	c.Ttl = runAndHttpPing(c)
}

func (c *ssrConfig) tcpPing() {
	c.Ttl = TcpPing(c.Host + ":" + strconv.Itoa(c.Port))
}

type ssrConfigSlice []*ssrConfig

func (s ssrConfigSlice) Len() int {
	return len(s)
}

func (s ssrConfigSlice) Less(i, j int) bool {
	return s[i].Ttl < s[j].Ttl
}

func (s ssrConfigSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func parseSSRUrl(url string) (*ssrConfig, error) {
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

func readSSR(data string) (map[string]*ssrConfig, error) {
	ssrs := strFilter(strings.Split(data, "\n"), func(s string) bool {
		return s != ""
	})

	cfgs := map[string]*ssrConfig{}
	for _, ssr := range ssrs {
		cfg, err := parseSSRUrl(ssr)
		if err != nil {
			log(err, ":", ssr)
			continue
		}
		cfgs[cfg.Remarks] = cfg
	}
	return cfgs, nil
}

func cfgsFromUrl(url string) (map[string]*ssrConfig, error) {
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

func checkCfg(master chan<- *ssrConfig, cfg *ssrConfig) {
	if checkAllMethod == "http" {
		cfg.httpPing()
	} else {
		cfg.tcpPing()
	}
	log(cfg.Remarks, cfg.Host, "ttl:", cfg.Ttl)
	master <- cfg
}

func goodWaysByCfgs(cfgs map[string]*ssrConfig, goodKeyWords []string, badKeyWords []string) ssrConfigSlice {
	var goodCfgs []*ssrConfig
	self := make(chan *ssrConfig)
	childNum := 0

	for _, cfg := range cfgs {
		log(cfg.Remarks)
	}

	for _, cfg := range cfgs {
		if (cfg.Method == "rc4-md5") {
			continue
		}
		if (len(badKeyWords) != 0 && anyStrsInStr(cfg.Remarks, badKeyWords)) ||
			(len(goodKeyWords) != 0 && !allStrsInStr(cfg.Remarks, goodKeyWords)) {
			log(cfg.Host, cfg.Remarks, "BAN")
			continue
		}

		go checkCfg(self, cfg)
		childNum += 1
	}

	for i := 0; i < childNum; i++ {
		//for _, cfg := range cfgs {
		cfg := <-self
		if cfg.Ttl > 0 {
			goodCfgs = append(goodCfgs, cfg)
		}
	}
	return goodCfgs
}

func goodWaysFromUrl(url string, goodKeyWords []string, badKeyWords []string) ssrConfigSlice {
	log("http get ssr config...")

	var (
		cfgMap map[string]*ssrConfig
		err    error
	)

	for {
		cfgMap, err = cfgsFromUrl(url)
		if err != nil {
			log("http get ssr config error:", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	log("SSR config length:", len(cfgMap))
	for {
		cfgs := goodWaysByCfgs(cfgMap, goodKeyWords, badKeyWords)
		if len(cfgs) == 0 {
			log("ssr configs all bad. again...")
			time.Sleep(5 * time.Second)
		} else {
			return cfgs
		}
	}
}

func bestWay(url string, goodKeyWords []string, badKeyWords []string) *ssrConfig {
	cfgs := goodWaysFromUrl(url, goodKeyWords, badKeyWords)

	ttlHostMap := map[int]*ssrConfig{}
	var ttls []int
	for _, cfg := range cfgs {
		ttlHostMap[cfg.Ttl] = cfg
		ttls = append(ttls, cfg.Ttl)
	}
	minTTL := minInt(ttls)
	best := ttlHostMap[minTTL]
	log("best addr:", best.Remarks, best.Host, best.Port, minTTL)

	return ttlHostMap[minTTL]
}

func goodWays(url string, goodKeyWords []string, badKeyWords []string) ssrConfigSlice {
	cfgs := goodWaysFromUrl(url, goodKeyWords, badKeyWords)

	sort.Sort(cfgs)

	log("good cfgs:")
	for _, cfg := range cfgs {
		log(cfg.Remarks, cfg.Ttl)
	}

	return cfgs
}
