//go:build integration

package flickr

import (
	"testing"
)

func TestFeed(t *testing.T) {
	client, err := NewPhotosClient()
	if err != nil {
		t.Fatalf("Unable to create client: %s", err)
	}

	photos, err := client.Feed("98269877@N00", 50)

	if err != nil {
		t.Fatalf("Error getting feed: %s", err)
	}

	if len(photos) < 40 {
		t.Fatalf("Less than 50 photos were returned: %d", len(photos))
	}

	if len(photos[0].Id) == 0 {
		t.Fatalf("First photo id is empty")
	}

	if len(photos[0].Owner) == 0 {
		t.Fatalf("First owner is empty")
	}
}
