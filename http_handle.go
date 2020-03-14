package main

import (
	"fmt"
	"net/http"
)

func httpIndex(w http.ResponseWriter, _ *http.Request) {
	ssrGateChan <- "node"
	node := <- httpServChan
	fmt.Fprintf(w, node)
}

func httpUdpate(w http.ResponseWriter, _ *http.Request) {
	ssrGateChan <- "update"
	fmt.Fprintf(w, "SSR GATE UPDATE")
}

