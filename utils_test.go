package main

import (
	"testing"
)

var (
	//s = "沪港专线  01 - 香港HKT11 1Gbps 0.1倍 Netflix HBO TVB"
	s = "香港 BGP 中继 2 - 香港HKT08打机 1Gbps 0.1倍 Netflix HBO TVB BAN"
	ss = []string{"HKT01", "0.1"}
)

func TestAnyStrsInStr(t *testing.T) {
	//goodKeyWords := []string{"0.1"}
	//badKeyWords := []string{""}

	//fmt.Printf("%", strings.Split(" a", " "))
	//}

	//log(anyStrsInStr(s, []string{}))
	if !anyStrsInStr(s, ss) {
		t.Fail()
	}
}

func TestAllStrsInStr(t *testing.T) {
	if allStrsInStr(s, ss) {
		t.Fail()
	}
}