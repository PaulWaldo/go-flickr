package flickr

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type Photo struct {
	Id          int
	Title       string
	Description string
	PostTS      int
	URL         string
	Format      string
	Views       int

	Secret string
	Server int
	Farm   int

	Username     string
	UserFullName string
	UserId       string
	UserLocation string
	UserIcon     string
}

type PhotoRaw struct {
	Photo struct {
		Id           string
		Secret       string
		Server       string
		Farm         int
		DateUploaded string

		IsFavorite   interface{}
		License      interface{}
		Safety_Level interface{}
		Rotation     interface{}

		OriginalSecret string
		OriginalFormat string

		Owner struct {
			NSID       string
			Username   string
			Realname   string
			Location   string
			IconServer string
			IconFarm   int
			Path_Alias string
		}

		Title       map[string]string
		Description map[string]string

		// who cares about this shit
		Visibility        interface{}
		Dates             interface{}
		Editability       interface{}
		PublicEditability interface{}
		Usage             interface{}
		Comments          interface{}
		Notes             interface{}
		People            interface{}
		Tags              interface{}

		URLs map[string][]map[string]string

		Media interface{}
		Stat  interface{}
	}
}

type SizeToSuffixMap map[int]string

// Suffix	Class	Longest edge (px)	Notes
// s	thumbnail	75	cropped square
// q	thumbnail	150	cropped square
// t	thumbnail	100
// m	small	240
// n	small	320
// w	small	400
// (none)	medium	500
// z	medium	640
// c	medium	800
// b	large	1024
// h	large	1600	has a unique secret; photo owner can restrict
// k	large	2048	has a unique secret; photo owner can restrict
// 3k	extra large	3072	has a unique secret; photo owner can restrict
// 4k	extra large	4096	has a unique secret; photo owner can restrict
// f	extra large	4096	has a unique secret; photo owner can restrict; only exists for 2:1 aspect ratio photos
// 5k	extra large	5120	has a unique secret; photo owner can restrict
// 6k	extra large	6144	has a unique secret; photo owner can restrict
// o	original	arbitrary	has a unique secret; photo owner can restrict; files have full EXIF data; files might not be rotated; files can use an arbitrary file extension

var PhotoSizes = SizeToSuffixMap{
	75:   "s",
	150:  "q",
	100:  "t",
	240:  "m",
	320:  "n",
	400:  "w",
	500:  "",
	640:  "z",
	800:  "c",
	1024: "b",
	1600: "h",
	2048: "k",
	3072: "3k",
	4096: "4k",
	5120: "5k",
	6144: "6k",
}

func closestSuffix(desiredSize int) (string, error) {
	// Sort PhotoSizes
	keys := make([]int, 0, len(PhotoSizes))
	for k := range PhotoSizes {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Search for key that is same or larger than desired
	for _, k := range keys {
		if k >= desiredSize {
			return PhotoSizes[k], nil
		}
	}
	return "", errors.New("unable to find matching size")
}

func (client *Client) GetURL(photo *RawFavedPhoto, desiredSize int) (string, error) {
	// # Typical usage
	// https://live.staticflickr.com/{server-id}/{id}_{secret}_{size-suffix}.jpg

	// # Unique URL format for 500px size
	// https://live.staticflickr.com/{server-id}/{id}_{secret}.jpg

	// # Originals might have a different file format extension
	// https://live.staticflickr.com/{server-id}/{id}_{o-secret}_o.{o-format}

	// # Example
	// #   server-id: 7372
	// #   photo-id: 12502775644
	// #   secret: acfd415fa7
	// #   size: w
	// https://live.staticflickr.com/7372/12502775644_acfd415fa7_w.jpg

	sizeSuffix, err := closestSuffix(desiredSize)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://live.staticflickr.com/%s/%s_%s_%s.jpg",
		photo.Server, photo.Id, photo.Secret, sizeSuffix)
	return url, nil
}

func (client *Client) GetPhoto(id int) (*Photo, error) {
	response, err := client.Request("photos.getInfo", Params{
		"photo_id": fmt.Sprintf("%d", id),
	})

	if err != nil {
		return nil, err
	}

	raw := &PhotoRaw{}
	err = Parse(response, raw)

	if err != nil {
		return nil, err
	}

	date, err := strconv.Atoi(raw.Photo.DateUploaded)

	if err != nil {
		return nil, err
	}

	server, err := strconv.Atoi(raw.Photo.Server)

	if err != nil {
		return nil, err
	}

	icon := fmt.Sprintf("https://farm%d.staticflickr.com/%s/buddyicons/%s_r.jpg", raw.Photo.Owner.IconFarm, raw.Photo.Owner.IconServer, raw.Photo.Owner.NSID)

	return &Photo{
		Id:          id,
		Title:       raw.Photo.Title["_content"],
		Description: raw.Photo.Description["_content"],
		PostTS:      date,
		Format:      raw.Photo.OriginalFormat,
		Secret:      raw.Photo.OriginalSecret,
		Server:      server,
		Farm:        raw.Photo.Farm,

		Username:     raw.Photo.Owner.Path_Alias,
		UserFullName: raw.Photo.Owner.Realname,
		UserId:       raw.Photo.Owner.NSID,
		UserLocation: raw.Photo.Owner.Location,
		UserIcon:     icon,
	}, nil
}
