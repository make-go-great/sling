package main

import (
	"fmt"
	"io"
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

	example(parent)
	exampleJSON(parent)
}

func example(parent *sling.Sling) {
	child, err := parent.Clone()
	if err != nil {
		fmt.Println(err)
		return
	}

	rsp, err := child.Response()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		fmt.Printf("Response status: %d\n", rsp.StatusCode)
		return
	}

	bytes, err := io.ReadAll(rsp.Body)
	if err != nil {
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
	if err := child.ResponseDecoder(&jsonRspDecoder).Receive(&u); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Result: %+v\n", u)
}
