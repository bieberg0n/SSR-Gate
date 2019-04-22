package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	logs "log"
	"os"
	"strings"
)

func log(args... interface{}) {
	logs.Println(args...)
}

func logb(arg interface{}) {
	data, err := json.MarshalIndent(arg, "", "  ")
	if err != nil {
		log(err)
		return
	}
	log(string(data))
}

func logj(args... interface{}) {
	for _, arg := range args {
		fmt.Printf("%+v\n", arg)
	}
}

func compareInt(arr []int, better func (int, int) bool) int {
	good := arr[0]
	for _, num := range arr[1:] {
		if better(good, num) {
			good = num
		}
	}
	return good
}

func maxInt(arr []int) int {
	return compareInt(arr, func(a int, b int) bool {return a < b})
}

func minInt(arr []int) int {
	return compareInt(arr, func(a int, b int) bool {return a > b})
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func intFilter(arr []int, isGood func (int) bool) []int {
	result := make([]int, 0)
	for _, i := range arr {
		if isGood(i) {
			result = append(result, i)
		}
	}
	return result

}

func strFilter(arr []string, isGood func (string) bool) []string {
	result := make([]string, 0)
	for _, i := range arr {
		if isGood(i) {
			result = append(result, i)
		}
	}
	return result
}

func strInclude(arr []string, str string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}

func pathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func b64decode(in string) (string, error) {
	in = strings.Replace(in, "_", "/", -1)
	in = strings.Replace(in, "-", "+", -1)
	buf := bytes.NewBufferString(in)
	for i := 0; i < buf.Len() % 4; i++ {
		buf.WriteString("=")
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(buf.String())
	if err != nil {
		return "", err
	}

	return string(decodeBytes), nil
}

func allStrsInStr(str string, strs []string) bool {
	for _, s := range strs {
		if s != "" && !strings.Contains(str, s) {
			return false
		}
	}
	return true
}

func anyStrsInStr(str string, strs []string) bool {
	for _, s := range strs {
		if s != "" && strings.Contains(str, s) {
			return true
		}
	}
	return false
}