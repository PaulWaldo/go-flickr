package main

import (
	"fmt"
	"os"

	"github.com/azer/go-flickr"
)

func main() {
	client, err := flickr.NewClient()
	if err != nil {
		fmt.Printf("Unable to create client: %s", err)
		os.Exit(1)
	}

	resp, err := client.Request("people.findByUsername", flickr.Params{
		"username": "azerbike",
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v", string(resp))
}
