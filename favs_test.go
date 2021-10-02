package flickr

import (
	"math"
	"testing"
)

func TestFavPagination(t *testing.T) {
	client := NewDefaultPaginatedClient("", "")

	favs, err := client.Favs("98269877@N00")
	if err != nil {
		t.Fatalf("Error getting Favs: %s", err)
	}

	if client.NumPages < 1 {
		t.Fatalf("expecting at least 1 page, got %d", client.NumPages)
	}
	if len(favs) < 1 {
		t.Fatalf("expecting at least 1 fav, got %d", len(favs))
	}

	// Only test 3 pages or the second to the last, whichever is smaller
	pagesToTest := math.Min(3, float64(client.NumPages))
	for i := 1; i <= int(pagesToTest); i++ {
		favs, err := client.NextPage()
		if err != nil {
			t.Fatalf("error getting Favs page %d of %d: %s", i, client.NumPages, err)
		}
		if client.Page != i+1 {
			t.Fatalf("expecting returned page to be %d but got %d", i+1, client.Page)
		}
		if len(favs) < 1 {
			t.Fatalf("expecting at least 1 fav, got %d", len(favs))
		}
	}

	// Read last page
	client.RequestPage = client.NumPages - 1
	_, err = client.NextPage()
	if err != nil {
		t.Fatalf("error getting Favs page %d of %d: %s", client.Page, client.NumPages, err)
	}
	if client.Page != client.NumPages {
		t.Fatalf("expecting returned page to be %d but got %d", client.RequestPage, client.Page)
	}

	// Next page should have a ErrPaginatorExhausted error
	_, err = client.NextPage()
	if err != ErrPaginatorExhausted {
		t.Fatalf("expecting error '%v' reading past end, but got %v", ErrPaginatorExhausted, err)
	}
}

func TestFavs(t *testing.T) {
	client := NewDefaultPaginatedClient("", "")

	favs, err := client.Favs("98269877@N00")
	if err != nil {
		t.Fatalf("error getting Favs: %s", err)
	}

	if client.NumPages < 1 {
		t.Fatalf("expecting at least 1 page, got %d", client.NumPages)
	}

	if client.Page != 1 {
		t.Fatalf("expecting at first page, got page %d", client.Page)
	}

	if client.Total < 1 {
		t.Fatalf("expecting at least 1 total fav, got %d", client.Total)
	}

	if len(favs) < 50 {
		t.Fatalf("less than 90 favorites were created: %d", len(favs))
	}

	if len(favs[0].Id) == 0 {
		t.Fatalf("first fav id is empty")
	}

	if len(favs[0].Owner) == 0 {
		t.Fatalf("first owner is empty")
	}
}
