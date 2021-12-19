//go:build integration
// +build integration

package flickr

import (
	"testing"
)

func TestFindUser(t *testing.T) {
	client, err := NewPhotosClient()
	if err != nil {
		t.Fatalf("Unable to create client: %s", err)
	}

	user, err := client.FindUser("azer")

	if err != nil {
		t.Fatalf("Error finding user: %s", err)
	}

	if user.Id != "98269877@N00" {
		t.Fatalf("Invalid user id")
	}

	if user.Name != "azer" {
		t.Fatalf("Invalid user name")
	}
}
