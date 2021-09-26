package flickr

import (
	"testing"
)

func TestFavs(t *testing.T) {
	client := NewDefaultPaginatedClient("", "")

	favs, err := client.Favs("98269877@N00")

	if err != nil {
		t.Fatalf("Error getting Favs: %s", err)
	}

	if favs.NumPages < 1 {
		t.Fatalf("Expecting at least 1 page, got %d", favs.NumPages)
	}
	if len(favs.Favs) < 50 {
		t.Fatalf("Less than 90 favorites were created: %d", len(favs.Favs))
	}

	if len(favs.Favs[0].Id) == 0 {
		t.Fatalf("First fav id is empty")
	}

	if len(favs.Favs[0].Owner) == 0 {
		t.Fatalf("First owner is empty")
	}
}
