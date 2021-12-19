//go:build integration
// +build integration

package flickr

import (
	"testing"
)

func TestAlbum(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Unable to create client: %s", err)
	}

	album, err := client.Album("72157617176794673")

	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(album.Id) == 0 {
		t.Fatalf("album id is empty")
	}

	if len(album.Photos) == 0 {
		t.Errorf("Expected to fing photos in album, found none")
	}

	if len(album.Photos[0].Id) == 0 {
		t.Fatalf("First photo id is empty")
	}
}
