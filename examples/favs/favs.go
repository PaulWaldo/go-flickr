package main

import (
	"fmt"
	"github.com/azer/go-flickr"
)

func main() {
	client := flickr.NewClient("", "")

	favs, err := client.Favs("98269877@N00")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d Favs. First title: %s", len(favs), favs[0].Title)
	}
}
