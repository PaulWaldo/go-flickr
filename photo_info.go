package flickr

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type PhotoInfo struct {
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

type PhotoInfoRaw struct {
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

func (client *Client) GetPhotoInfo(id int) (*PhotoInfo, error) {
	response, err := client.Request("photos.getInfo", Params{
		"photo_id": fmt.Sprintf("%d", id),
	})

	if err != nil {
		return nil, err
	}

	raw := &PhotoInfoRaw{}
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

	return &PhotoInfo{
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

var ErrMinSizeNotAvailable = errors.New("minimum size not available")

func (ps *PhotoSizes) ClosestWidthUrl(desired, minAccepted int) (string, error) {
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

	difAllowed := desired - minAccepted
	for _, v := range keys {
		if v >= desired || v >= desired-difAllowed {
			return widthUrlMap[v], nil
		}
	}

	// largestAvailable := keys[len(keys)-1]
	return "", ErrMinSizeNotAvailable
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
