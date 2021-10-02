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
