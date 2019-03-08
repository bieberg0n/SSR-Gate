package main

import (
	"testing"
)

//func TestReadSSR (t *testing.T) {
//	cfgs, err := readSSR(ssrData)
//	if err != nil {
//		t.Log(err)
//		return
//	}
//	for _, cfg := range cfgs {
//		logj(cfg)
//	}
//}

func TestParseSSRUrl (t *testing.T) {
	cfg, err := parseSSRUrl("ssr://NDUuNjIuMjM4LjE0Nzo1NjA1OmF1dGhfc2hhMV92NDpjaGFjaGEyMDp0bHMxLjJfdGlja2V0X2F1dGg6Wkc5MVlpNXBieTl6YzNwb1puZ3ZLalUyTURVLz9yZW1hcmtzPTVweXM1WVdONkxTNTZMU201WS0zNXAybDZJZXFPbVJ2ZFdJdWFXOHZjM042YUdaNEx3")
	logs(err)
	logb(cfg)
}