package main

import (
	"fmt"
	"os"

	"github.com/azer/go-flickr"
)

func main() {
	client, err := flickr.NewPhotosClient()
	if err != nil {
		fmt.Printf("Unable to create client: %s", err)
		os.Exit(1)
	}

	favs, err := client.Favs("98269877@N00")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Favs have %d pages of %d each\n", favs.Pages, favs.PerPage)
		fmt.Printf("First title of page %d: %s\n", favs.Page, favs.Photos[0].Title)
	}

	if favs.Pages > 1 {
		favs, err = client.NextPage()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("First title of page %d: %s\n", favs.Page, favs.Photos[0].Title)
		}
	}
}
