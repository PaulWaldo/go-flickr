package flickr

import (
	"testing"
)

func TestGetPhoto(t *testing.T) {
	client := NewClient("", "")

	photo, err := client.GetPhoto(15691826511)

	if err != nil {
		t.Fatalf("Error getting phptp: %s", err)
	}

	if photo.Id != 15691826511 {
		t.Fatalf("Invalid photo")
	}

	if photo.Title != "driving to cusco" {
		t.Fatalf("Invalid photo title")
	}

	if photo.Username != "azer" {
		t.Fatalf("Invalid user")
	}

	// if photo.UserIcon != "https://farm3.staticflickr.com/2933/buddyicons/98269877@N00_r.jpg" {
	// 	t.Error("Invalid user icon")
	// }
}
