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

	following, err := client.Following("98269877@N00")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d Following. First user: %s", len(following), following[0].Title)
	}
}
