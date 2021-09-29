package flickr

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Fav struct {
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
		Photo []Fav
	}
}

type paginationState struct {
	userId string
}

func (client *PaginatedClient) Favs(userId string) ([]Fav, error) {
	response, err := client.Request("favorites.getPublicList", Params{
		"user_id": userId,
	})
	if err != nil {
		return nil, err
	}

	raw := &favsRaw{}
	err = ParsePaginated(response, &raw.Photos.PaginatedResult, raw)
	if err != nil {
		return nil, err
	}

	client.PaginationState = paginationState{
		userId: userId,
	}
	client.NumPages = raw.Photos.NumPages
	client.Total = raw.Photos.Total
	client.Page = raw.Photos.Page
	client.NumPerPage = raw.Photos.NumPerPage

	return raw.Photos.Photo, nil
}

func (client *PaginatedClient) NextPage() ([]Fav, error) {
	if client.RequestPage > client.NumPages {
		return nil, ErrPaginatorExhausted
	}

	state, ok := client.PaginationState.(paginationState)
	if !ok {
		return nil, fmt.Errorf("PaginationState is wrong type: %v", state)
	}

	favs, err := client.Favs(state.userId)
	if err != nil {
		return nil, fmt.Errorf("error getting Favs: %s", err)
	}
	return favs, nil
}

func divMod(dvdn, dvsr int) (q, r int) {
	r = dvdn
	for r >= dvsr {
		q += 1
		r = r - dvsr
	}
	return
}

const allRightReserved = "0"

func (client *Client) RandomFav(userId string) (Fav, error) {
	const pageSize = 100
	response, err := client.Request("favorites.getPublicList", Params{
		"user_id": userId, "per_page": strconv.Itoa(pageSize), "extras": "license",
	})
	if err != nil {
		return Fav{}, err
	}

	raw := &favsRaw{}
	err = Parse(response, raw)

	if err != nil {
		return Fav{}, err
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
			return Fav{}, err
		}

		raw = &favsRaw{}
		err = Parse(response, raw)

		if err != nil {
			return Fav{}, err
		}

		isReserved = raw.Photos.Photo[offset].License == allRightReserved
	}
	return raw.Photos.Photo[offset], nil
}
