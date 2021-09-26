package main

import (
	"fmt"
	"github.com/azer/go-flickr"
)

func main() {
	client := flickr.NewDefaultPaginatedClient("", "")

	favs, err := client.Favs("98269877@N00")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d Favs. First title: %s", len(favs.Favs), favs.Favs[0].Title)
	}
}
