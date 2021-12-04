package flickr

import (
	"testing"
)

func TestGetPhoto(t *testing.T) {
	client, err := NewPhotosClient()
	if err != nil {
		t.Fatalf("Unable to create client: %s", err)
	}

	const photoId int = 51387563129

	photo, err := client.GetPhotoInfo(photoId)

	if err != nil {
		t.Fatalf("Error getting photo: %s", err)
	}

	if photo.Id != photoId {
		t.Fatalf("Invalid photo")
	}

	if photo.Title != "Sunflowers forever" {
		t.Fatalf("Invalid photo title")
	}

	if photo.Username != "geekneck" {
		t.Fatalf("Expecting user geekneck, got %s", photo.Username)
	}
}

func TestGetPhotoGetSizes(t *testing.T) {
	client, err := NewPhotosClient()
	if err != nil {
		t.Fatalf("Unable to create client: %s", err)
	}

	sizes, err := client.GetPhotoSizes(567229075)
	if err != nil {
		t.Fatalf("Error getting photo sizes: %s", err)
	}
	if sizes.CanBlog != 0 {
		t.Fatalf("Expecting CanBlog to be 0, but got %d", sizes.CanBlog)
	}
	if sizes.CanDownload != 1 {
		t.Fatalf("Expecting CanDownload to be 1, but got %d", sizes.CanDownload)
	}
	if sizes.CanPrint != 0 {
		t.Fatalf("Expecting CanPrint to be 0, but got %d", sizes.CanPrint)
	}
	square := sizes.Sizes[0]
	if square.Height != 75 {
		t.Fatalf("Expecting Height to be 75, but got %d", square.Height)
	}
	if square.Width != 75 {
		t.Fatalf("Expecting Width to be 75, but got %d", square.Width)
	}
	if square.Label != "Square" {
		t.Fatalf("Expecting Label to be Square, but got %s", square.Label)
	}
	expectedSource := "https://live.staticflickr.com/1103/567229075_2cf8456f01_s.jpg"
	if square.Source != expectedSource {
		t.Fatalf("Expecting Source to be %s, but got %s", expectedSource, square.Source)
	}
	expectedUrl := "https://www.flickr.com/photos/stewart/567229075/sizes/sq/"
	if square.URL != expectedUrl {
		t.Fatalf("Expecting URL to be %s, but got %s", expectedUrl, square.URL)
	}
	if square.Media != "photo" {
		t.Fatalf("Expecting Media to be photo, but got %s", square.Media)
	}
}

const smallUrl = "smallUrl"
const mediumUrl = "mediumUrl"
const largeUrl = "largeUrl"

var photoSizeInfo = &PhotoSizes{
	CanBlog:     1,
	CanPrint:    0,
	CanDownload: 1,
	Sizes: []PhotoSize{
		{
			Height: 101,
			Width:  100,
			Label:  "Small",
			Media:  "photo",
			Source: smallUrl,
			URL:    "",
		},
		{
			Height: 201,
			Width:  200,
			Label:  "Medium",
			Media:  "photo",
			Source: mediumUrl,
			URL:    "",
		},
		{
			Height: 301,
			Width:  300,
			Label:  "Large",
			Media:  "photo",
			Source: largeUrl,
			URL:    "",
		},
	},
}

func TestClosestWidthUrl(t *testing.T) {
	cases := []struct {
		desired     int
		minAccepted int
		expected    string
		errExpected bool
	}{
		{1, 1, smallUrl, false},
		{99, 99, smallUrl, false},
		{100, 100, smallUrl, false},
		{101, 101, mediumUrl, false},
		{99999, 99999, largeUrl, true},
		{101, 100, smallUrl, false},
		{302, 302, "", true},
	}

	for _, c := range cases {
		actual, err := photoSizeInfo.ClosestWidthUrl(c.desired, c.minAccepted)
		if err != nil {
			if !c.errExpected {
				t.Fatalf("Expecting error %v, but got %v", c.errExpected, err)
			} else {
				continue
			}
		} else {
			if c.errExpected {
				t.Fatal("Expecting error but got none")
			}
		}
		if actual != c.expected {
			t.Fatalf("For requested size %d, minAccepted %d, Expecting url to be '%s', but got '%s'", c.desired, c.minAccepted, c.expected, actual)
		}
	}
}
