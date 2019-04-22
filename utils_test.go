package main

import (
	"testing"
)

var (
	s = "沪港专线  01 - 香港HKT11 1Gbps 0.1倍 Netflix HBO TVB"
	ss = []string{"HKT01", "0.1"}
)

func TestAnyStrsInStr(t *testing.T) {
	if !anyStrsInStr(s, ss) {
		t.Fail()
	}
}

func TestAllStrsInStr(t *testing.T) {
	if allStrsInStr(s, ss) {
		t.Fail()
	}
}