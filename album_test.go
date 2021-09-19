package flickr

import (
	"fmt"
	"testing"
)

func TestAlbum(t *testing.T) {
	client := NewClient("", "")

	album, err := client.Album("72157617176794673")

	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}

	fmt.Println(album.Id)
	fmt.Println(len(album.Photos))

	if len(album.Id) == 0 {
		t.Errorf("album id is empty")
	}

	if len(album.Photos) == 0 {
		t.Errorf("Expected to fing photos in album, found none")
	}

	if len(album.Photos[0].Id) == 0 {
		t.Errorf("First photo id is empty")
	}

	fmt.Println(album.Photos[0].URLs()["large"])
}
