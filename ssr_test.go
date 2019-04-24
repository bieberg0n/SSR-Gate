package main

import (
	"testing"
)

func TestParseSSRUrl (t *testing.T) {
	cfg, err := parseSSRUrl("ssr://MjYwNDphODgwOmNhZDpkMDo6MjlmOmYwMDE6NjA0MDA6b3JpZ2luOmFlcy0yNTYtY2ZiOnBsYWluOmQwbHBZbUpJLz9vYmZzcGFyYW09JnByb3RvcGFyYW09JnJlbWFya3M9NVlxZzVvdV81YVNuSUdsd2RqWWdNMGRpY0hNZ01lV0FqUSZncm91cD1jM055WTJ4dmRXUQ")
	log(err)
	logb(cfg)
}