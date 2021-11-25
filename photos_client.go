package flickr

import (
	"encoding/json"
	"strconv"
)

// Bool allows 0/1 to also become boolean.
type Bool bool

// Courtesy of https://stackoverflow.com/a/56832346/1290460
func (bit *Bool) UnmarshalJSON(b []byte) error {
	txt := string(b)
	*bit = Bool(txt == "1" || txt == "true")
	return nil
}

/*
	Paginator represent the paged data that Flickr returns when presenting photos.
	When photos are requested, they arrive in batches (pages) with an
	enclosing structure:
	<photos page="2" pages="89" perpage="10" total="881">
		<photo id="2636" owner="47058503995@N01"
			secret="a123456" server="2" title="test_04"
			ispublic="1" isfriend="0" isfamily="0" />
		<photo id="2635" owner="47058503995@N01"
			secret="b123456" server="2" title="test_03"
			ispublic="0" isfriend="1" isfamily="1" />
		<photo id="2633" owner="47058503995@N01"
			secret="c123456" server="2" title="test_01"
			ispublic="1" isfriend="0" isfamily="0" />
		<photo id="2610" owner="12037949754@N01"
			secret="d123456" server="2" title="00_tall"
			ispublic="1" isfriend="0" isfamily="0" />
	</photos>
*/

// Server side data structure

type photoListRaw struct {
	Photos paginatedResult
}

type paginatedResult struct {
	Page       int `json:"page"`
	NumPages   int `json:"pages"`
	NumPerPage int `json:"perpage"`
	Total      int `json:"total"`
	Photo      []PhotoListItem
}

// Our data structure

type Paginator struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	PerPage int `json:"perpage"`
	Total   int `json:"total"`
}

type PaginationParams struct {
	PerPage int
	Page    int
}

type PhotoListItem struct {
	ID       string `json:"id"`
	Owner    string
	Secret   string
	Server   string
	Title    string
	IsPublic Bool
	IsFriend Bool
	IsFamily Bool
}

type PhotoList struct {
	Paginator
	Photos []PhotoListItem
}

type context struct {
	method string
	params Params
}

type PhotosClient struct {
	*Client
	*PaginationParams
	context
}

func NewPhotosClientFull(apiKey string, envFileName string, paginationParams PaginationParams) *PhotosClient {
	return &PhotosClient{
		Client:           NewClient(apiKey, envFileName),
		PaginationParams: &paginationParams,
	}
}

func NewPhotosClient() *PhotosClient {
	return &PhotosClient{
		Client:           NewClient("", ""),
		PaginationParams: &PaginationParams{PerPage: 100, Page: 1},
	}
}
func convert(r photoListRaw) *PhotoList {
	x := &PhotoList{}
	x.Page = r.Photos.Page
	x.Pages = r.Photos.NumPages
	x.PerPage = r.Photos.NumPerPage
	x.Total = r.Photos.Total
	x.Photos = r.Photos.Photo
	return x
}

func (c *PhotosClient) Request(method string, params Params) (*PhotoList, error) {
	params["page"] = strconv.Itoa(c.Page)
	params["per_page"] = strconv.Itoa(c.PerPage)

	b, err := c.Client.Request(method, params)
	if err != nil {
		return nil, err
	}

	v := &photoListRaw{}
	err = json.Unmarshal(b, v)
	if err != nil {
		return nil, err
	}

	// Save current context for a potential call to NextPage
	c.params = params
	c.method = method

	return convert(*v), nil
}

func (c *PhotosClient) NextPage() (*PhotoList, error) {
	c.PaginationParams.Page = c.Page + 1
	return c.Request(c.method, c.params)
}
