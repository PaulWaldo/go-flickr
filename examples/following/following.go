package main

import (
	"fmt"

	"github.com/azer/go-flickr"
)

func main() {
	client := flickr.NewClient("", "")

	following, err := client.Following("98269877@N00")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d Following. First user: %s", len(following), following[0].Title)
	}
}
