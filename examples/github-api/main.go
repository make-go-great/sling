package main

import (
	"net/http"

	"github.com/haunt98/sling"
)

func main() {
	s := sling.New(http.DefaultClient)
	s.Get("https://api.github.com/haunt98").
		AddHeader("Accept", "application/vnd.github.v3+json")
}
