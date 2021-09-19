package main

import (
	"fmt"

	"github.com/azer/go-flickr"
)

func main() {
	client := flickr.NewClient("", "")

	resp, err := client.Request("people.findByUsername", flickr.Params{
		"username": "azerbike",
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v", string(resp))
}
