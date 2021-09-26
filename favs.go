package flickr

import (
	"math/rand"
	"strconv"
)

// type Fav struct {
// 	Id        string
// 	Title     string
// 	Owner     string
// 	FavedBy   string
// 	DateFaved string
// 	Farm      int
// 	Secret    string
// 	Server    string
// }

type rawFavedPhoto struct {
	DateFaved string `json:"date_faved"`
	Farm      int
	Id        string
	IsFamily  int
	IsFriend  int
	License   string
	Owner     string
	Secret    string
	Server    string
	Title     string
	IsPublic  int
}

type favsRaw struct {
	Photos struct {
		PaginatedResult
		Photo []rawFavedPhoto
	}
	// Stat string
}

type FavsResponse struct {
	PaginatedResult
	Favs []rawFavedPhoto
}

func responseFromRaw(raw *favsRaw) *FavsResponse {
	resp := &FavsResponse{}
	resp.Favs = raw.Photos.Photo
	resp.PaginatedResult = raw.Photos.PaginatedResult
	return resp
}
// type favsRaw2 struct {
// 	Photos struct {
// 		Page    int
// 		Pages   int
// 		PerPage int
// 		Photo   []rawFavedPhoto
// 		Total   string
// 	}
// 	Stat string
// }

func (client *PaginatedClient) Favs(userId string) (*FavsResponse, error) {
	response, err := client.Request("favorites.getPublicList", Params{
		"user_id": userId,
	})
	if err != nil {
		return nil, err
	}

	raw := &favsRaw{}
	err = ParsePaginated(response, &raw.Photos.PaginatedResult, raw)
	// if err != nil {
	// 	raw := &FavsRaw2{}
	// 	err = Parse(response, raw)
	// }
	if err != nil {
		return nil, err
	}

	// favs := []Fav{}

	// for _, photo := range raw.Photos.Photo {
	// 	favs = append(favs, Fav{
	// 		Id:        photo.Id,
	// 		Title:     photo.Title,
	// 		Owner:     photo.Owner,
	// 		FavedBy:   userId,
	// 		DateFaved: photo.Date_Faved,
	// 		Farm:      photo.Farm,
	// 		Secret:    photo.Secret,
	// 		Server:    photo.Server,
	// 	})
	// }

	resp := responseFromRaw(raw)

	return resp, nil
}

func divMod(dvdn, dvsr int) (q, r int) {
	r = dvdn
	for r >= dvsr {
		q += 1
		r = r - dvsr
	}
	return
}

const AllRightReserved = "0"

func (client *Client) RandomFav(userId string) (rawFavedPhoto, error) {
	const pageSize = 100
	response, err := client.Request("favorites.getPublicList", Params{
		"user_id": userId, "per_page": strconv.Itoa(pageSize), "extras": "license",
	})
	if err != nil {
		return rawFavedPhoto{}, err
	}

	raw := &favsRaw{}
	err = Parse(response, raw)

	// if err != nil {
	// 	raw := &FavsRaw2{}
	// 	err = Parse(response, raw)
	// }

	if err != nil {
		return rawFavedPhoto{}, err
	}

	// Loop through random Favs, looking for ones which are not restricted, i.e.
	// license value != 0 ("All Rights Reserved")
	var page, offset int
	isReserved := true
	for isReserved {
		photoNum := rand.Intn(raw.Photos.Total)
		page, offset = divMod(photoNum, pageSize)
		page += 1 // Account for API pages starting at 1
		response, err = client.Request("favorites.getPublicList", Params{
			"user_id":  userId,
			"per_page": strconv.Itoa(pageSize),
			"page":     strconv.Itoa(page),
			"extras":   "license",
		})
		if err != nil {
			return rawFavedPhoto{}, err
		}

		raw = &favsRaw{}
		err = Parse(response, raw)

		// if err != nil {
		// 	raw := &favsRaw2{}
		// 	err = Parse(response, raw)
		// }

		if err != nil {
			return rawFavedPhoto{}, err
		}

		isReserved = raw.Photos.Photo[offset].License == AllRightReserved
	}
	return raw.Photos.Photo[offset], nil
}
