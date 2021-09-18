package flickr

import (
	"math/rand"
	"strconv"
)

type Fav struct {
	Id        string
	Title     string
	Owner     string
	FavedBy   string
	DateFaved string
	Farm      int
	Secret    string
	Server    string
}

type RawFavedPhoto struct {
	Date_Faved string
	Farm       int
	Id         string
	IsFamily   int
	IsFriend   int
	IsPublic   int
	Owner      string
	Secret     string
	Server     string
	Title      string
}

type FavsRaw struct {
	Photos struct {
		Page    int
		Pages   int
		PerPage int
		Total   int
		Photo   []RawFavedPhoto
	}
	Stat string
}

type FavsRaw2 struct {
	Photos struct {
		Page    int
		Pages   int
		PerPage int
		Photo   []RawFavedPhoto 
		Total string
	}
	Stat string
}

func (client *Client) Favs(userId string) ([]Fav, error) {
	response, err := client.Request("favorites.getPublicList", Params{
		"user_id": userId,
	})
	if err != nil {
		return nil, err
	}

	raw := &FavsRaw{}
	err = Parse(response, raw)

	if err != nil {
		raw := &FavsRaw2{}
		err = Parse(response, raw)
	}

	if err != nil {
		return nil, err
	}

	favs := []Fav{}

	for _, photo := range raw.Photos.Photo {
		favs = append(favs, Fav{
			Id:        photo.Id,
			Title:     photo.Title,
			Owner:     photo.Owner,
			FavedBy:   userId,
			DateFaved: photo.Date_Faved,
			Farm:      photo.Farm,
			Secret:    photo.Secret,
			Server:    photo.Server,
		})
	}

	return favs, nil
}

func DivMod(dvdn, dvsr int) (q, r int) {
	r = dvdn
	for r >= dvsr {
		q += 1
		r = r - dvsr
	}
	return
}

func (client *Client) RandomFav(userId string) (RawFavedPhoto, error) {
	const pageSize = 100
	response, err := client.Request("favorites.getPublicList", Params{
		"user_id": userId, "per_page": strconv.Itoa(pageSize),
	})
	if err != nil {
		return RawFavedPhoto{}, err
	}

	raw := &FavsRaw{}
	err = Parse(response, raw)

	if err != nil {
		raw := &FavsRaw2{}
		err = Parse(response, raw)
	}

	if err != nil {
		return RawFavedPhoto{}, err
	}

	photoNum := rand.Intn(raw.Photos.Total)
	page, offset := DivMod(photoNum, pageSize)
	page += 1 // Account for API pages starting at 1
	response, err = client.Request("favorites.getPublicList", Params{
		"user_id": userId, "per_page": strconv.Itoa(pageSize), "page": strconv.Itoa(page),
	})
	if err != nil {
		return RawFavedPhoto{}, err
	}

	raw = &FavsRaw{}
	err = Parse(response, raw)

	if err != nil {
		raw := &FavsRaw2{}
		err = Parse(response, raw)
	}

	if err != nil {
		return RawFavedPhoto{}, err
	}
	return raw.Photos.Photo[offset], nil
}
