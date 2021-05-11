package main

import (
	"fmt"
	"net/http"

	"github.com/haunt98/sling"
	slinghttp "github.com/haunt98/sling/http"
)

type user struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

func main() {
	parent := sling.New(http.DefaultClient).
		Get("https://api.github.com/users/haunt98").
		AddHeader("Accept", "application/vnd.github.v3+json")

	exampleRaw(parent)
	exampleJSON(parent)
}

func exampleRaw(parent *sling.Sling) {
	child, err := parent.Clone()
	if err != nil {
		fmt.Println(err)
		return
	}

	var s string
	var rawRspDecoder slinghttp.RawResponseDecoder
	if err := child.ResponseDecoder(&rawRspDecoder).Receive(&s); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Result: %s\n", s)
}

func exampleJSON(parent *sling.Sling) {
	child, err := parent.Clone()
	if err != nil {
		fmt.Println(err)
		return
	}

	var u user
	var jsonRspDecoder slinghttp.JSONResponseDecoder
	if err := child.ResponseDecoder(&jsonRspDecoder).Receive(&u); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Result: %+v\n", u)
}
