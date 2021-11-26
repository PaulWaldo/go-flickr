package main

import (
	"fmt"
	"os"

	"github.com/azer/go-flickr"
)

func main() {
	client,err := flickr.NewClient()
	if err != nil {
		fmt.Printf("Unable to create client: %s", err)
		os.Exit(1)
	}

	user, err := client.FindUser("azer")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s: %s", user.Id, user.Name)
	}
}
