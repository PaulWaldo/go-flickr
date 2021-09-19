package main

import (
	"fmt"

	"github.com/azer/go-flickr"
)

func main() {
	client := flickr.NewClient("", "")

	user, err := client.FindUser("azer")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s: %s", user.Id, user.Name)
	}
}
