package flickr

import (
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

	// sizeSuffix, err := closestSuffix(desiredSize)
	// if err != nil {
	// 	return "", err
	// }
	url := fmt.Sprintf("https://live.staticflickr.com/%s/%s_%s_%s_d.jpg",
		photo.Server, photo.Id, photo.Secret, "w" /*sizeSuffix*/)
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

type PhotoSize struct {
	Height int    `json:"height"`
	Label  string `json:"label"`
	Media  string `json:"media"`
	Source string `json:"source"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type PhotoSizes struct {
	CanBlog     int `json:"canblog"`
	CanPrint    int `json:"canprint"`
	CanDownload int `json:"candownload"`
	Sizes       []PhotoSize
}

func (ps *PhotoSizes) ClosestWidthUrl(desired int) (string, error) {
	widthUrlMap := make(map[int]string)
	for _, s := range ps.Sizes {
		widthUrlMap[s.Width] = s.Source
	}

	// Sort keys
	keys := make([]int, 0, len(ps.Sizes))
	for k := range widthUrlMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, v := range keys {
		if v >= desired {
			return widthUrlMap[v], nil
		}
	}

	largestAvailable := keys[len(keys)-1]
	return "", fmt.Errorf("largest size is %d, but %d requested", largestAvailable, desired)
}

type rawPhotoSizes struct {
	Sizes struct {
		CanBlog     int `json:"canblog"`
		CanPrint    int `json:"canprint"`
		CanDownload int `json:"candownload"`
		Size        []PhotoSize
	}
}

func (r *rawPhotoSizes) convert() *PhotoSizes {
	real := &PhotoSizes{}
	real.CanBlog = r.Sizes.CanBlog
	real.CanPrint = r.Sizes.CanPrint
	real.CanDownload = r.Sizes.CanDownload
	real.Sizes = r.Sizes.Size
	return real
}

func (client *Client) GetPhotoSizes(id int) (*PhotoSizes, error) {
	response, err := client.Request("photos.getSizes", Params{
		"photo_id": fmt.Sprintf("%d", id),
	})

	if err != nil {
		return nil, err
	}

	raw := &rawPhotoSizes{}
	err = Parse(response, raw)
	if err != nil {
		return nil, err
	}
	return raw.convert(), nil

}
