package flickr

import (
	"fmt"
	"testing"
)

func TestFavs(t *testing.T) {
	client := NewClient("", "")

	favs, err := client.Favs("98269877@N00")

	if err != nil {
		t.Error(err)
	}

	if len(favs) < 50 {
		t.Error(fmt.Printf("Less than 90 favorites were created: %d", len(favs)))
	}

	if len(favs[0].Id) == 0 {
		t.Error("First fav id is empty")
	}

	if len(favs[0].Owner) == 0 {
		t.Error("First owner is empty")
	}
}
