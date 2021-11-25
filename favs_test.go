package flickr

import (
	"testing"
)

func TestFavsPage1(t *testing.T) {
	client := NewPhotosClient()
	favs, err := client.Favs("98269877@N00")
	if err != nil {
		t.Fatalf("Error getting Favs: %s", err)
	}

	if favs.Page != 1 {
		t.Fatalf("Expecting page to be 1, but got %d", favs.Page)
	}
	if favs.Pages < 1 {
		t.Fatalf("Expecting at least 1 page, but got %d", favs.Pages)
	}
	if favs.PerPage != 100 {
		t.Fatalf("Expecting 100 photos per page, but got %d", favs.PerPage)

	}
	if len(favs.Photos) < 1 {
		t.Fatalf("Expecting 1 or more photos, got %d", len(favs.Photos))
	}
}

func TestFavsPage2(t *testing.T) {
	client := NewPhotosClient()
	client.PaginationParams.Page = 2
	favs, err := client.Favs("98269877@N00")
	if err != nil {
		t.Fatalf("Error getting Favs: %s", err)
	}

	if favs.Page != 2 {
		t.Fatalf("Expecting page to be 2, but got %d", favs.Page)
	}
	if len(favs.Photos) < 1 {
		t.Fatalf("Expecting 1 or more photos, got %d", len(favs.Photos))
	}
}

func TestPaginatorExhausted(t *testing.T) {
	client := NewPhotosClient()
	favs, err := client.Favs("98269877@N00")
	if err != nil {
		t.Fatalf("Error getting Favs: %s", err)
	}

	// Attempt to read past the last page
	client.PaginationParams.Page = favs.Pages
	_, err = client.NextPage()
	if err != ErrPaginatorExhausted {
		t.Fatalf("Expecting ErrPaginatorExhausted but got %v", err)
	}

}

func TestNextPage(t *testing.T) {
	client := NewPhotosClient()
	favs, err := client.Favs("98269877@N00")
	if err != nil {
		t.Fatalf("Error getting Favs: %s", err)
	}

	if favs.Page != 1 {
		t.Fatalf("Expecting page to be 1, but got %d", favs.Page)
	}
	if favs.Pages < 1 {
		t.Fatalf("Expecting for least 1 page, but got %d", favs.Pages)
	}
	totalPages := favs.Pages

	favs, err = client.NextPage()
	if err != nil {
		t.Fatalf("Error getting Favs page 2: %s", err)
	}

	if favs.Page != 2 {
		t.Fatalf("Expecting page to be 2, but got %d", favs.Page)
	}
	if favs.Pages != totalPages {
		t.Fatalf("Expecting total pages to be %d, but got %d", totalPages, favs.Pages)
	}
}

