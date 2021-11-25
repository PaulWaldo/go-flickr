package main

import (
	"fmt"

	"github.com/azer/go-flickr"
)

func main() {
	client := flickr.NewPhotosClient()
	favs, err := client.Favs("98269877@N00")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Favs have %d pages of %d each\n", favs.Pages, favs.PerPage)
		fmt.Printf("First title of page %d: %s\n", favs.Page, favs.Photos[0].Title)
	}

	if favs.Pages > 1 {
		favs, err = client.NextPage()
		fmt.Printf("favs=%+v\n", favs)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("First title of page %d: %s\n", favs.Page, favs.Photos[0].Title)
		}
	}
}
