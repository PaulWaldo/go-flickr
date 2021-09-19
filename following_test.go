package flickr

import (
	"testing"
)

func TestFollowing(t *testing.T) {
	client := NewClient("", "")

	following, err := client.Following("98269877@N00")

	if err != nil {
		t.Fatalf("Erro getting followers: %s", err)
	}

	if len(following) < 90 {
		t.Fatalf("Less than 400 following were found")
	}

	if len(following[0].Id) == 0 {
		t.Fatalf("First user id is empty")
	}

	if len(following[0].Title) == 0 {
		t.Fatalf("First user title is empty")
	}
}
