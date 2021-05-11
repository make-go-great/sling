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
		BaseURL("https://api.github.com").
		SetHeader("Accept", "application/vnd.github.v3+json")

	exampleRaw(parent)
	exampleJSON(parent)
}

func exampleRaw(parent *sling.Sling) {
	child, err := parent.Clone()
	if err != nil {
		fmt.Println(err)
		return
	}

	var bytes []byte
	var rawRspDecoder slinghttp.RawResponseDecoder
	if err := child.Get("/repos/haunt98/sling").ResponseDecoder(&rawRspDecoder).Receive(&bytes); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Result: %s\n", string(bytes))
}

func exampleJSON(parent *sling.Sling) {
	child, err := parent.Clone()
	if err != nil {
		fmt.Println(err)
		return
	}

	var u user
	var jsonRspDecoder slinghttp.JSONResponseDecoder
	if err := child.Get("/users/haunt98").ResponseDecoder(&jsonRspDecoder).Receive(&u); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Result: %+v\n", u)
}
