package main

import "testing"

func TestHttpPing(t *testing.T) {
	t.Log(HttpPing(1081))
}
