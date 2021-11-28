package flickr

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Bool allows 0/1 to also become boolean.
type Bool bool

// UnmarshalJSON unmarshals a Bool from JSON
// Courtesy of https://stackoverflow.com/a/56832346/1290460
func (bit *Bool) UnmarshalJSON(b []byte) error {
	txt := string(b)
	*bit = Bool(txt == "1" || txt == "true")
	return nil
}

// Server side data structures

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

// ErrPaginatorExhausted indicates that an attempt was made go beyond the last page of a paginated result
var ErrPaginatorExhausted = errors.New("attempt to read past last page of data")

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
type Paginator struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	PerPage int `json:"perpage"`
	Total   int `json:"total"`
}

// PaginationParams request Flickr to return results in chunk sizes
type PaginationParams struct {
	// PerPage is how many results to return at one time
	PerPage int
	// Page is the page that is being requested
	Page int
}

// PhotoListItem is the information about a specific photo
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

// PhotoList is the combined response from Flickr containg the paging data as well as the list of photos
type PhotoList struct {
	Paginator
	Photos []PhotoListItem
}

type context struct {
	method     string
	params     Params
	totalPages int
}

// PhotosClient is a client used to get Paginated data from Flickr
type PhotosClient struct {
	*Client
	*PaginationParams
	context
}

// NewPhotosClientEnvVar creates a paginated client that can access the Flickr API,
// by searching for an environment variable named by ApiKeyEnvVar
func NewPhotosClientEnvVar(string, paginationParams PaginationParams) (*PhotosClient, error) {
	client, err := NewClientEnvVar()
	if err != nil {
		return nil, err
	}
	return &PhotosClient{
		Client:           client,
		PaginationParams: &paginationParams,
	}, nil
}

// NewPhotosClientEnvFile creates a paginated client that can access the Flickr API,
// attempting fo fetch the API Key first from an environment file specified in envFileName
// then from the file ./.env
func NewPhotosClientEnvFile(envFileName string, paginationParams PaginationParams) (*PhotosClient, error) {
	client, err := NewClientEnvFile(envFileName)
	if err != nil {
		return nil, err
	}
	return &PhotosClient{
		Client:           client,
		PaginationParams: &paginationParams,
	}, nil
}

// NewClient creates a paginated client that can access the Flickr API,
// attempting fo fetch the API Key from the file ./.env
func NewPhotosClient() (*PhotosClient, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	return &PhotosClient{
		Client:           client,
		PaginationParams: &PaginationParams{PerPage: 100, Page: 1},
	}, nil
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

// Request uses the client to call the Flickr API.
// method is the name of the API to call.
// params are extra values used to modify the request
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

	conv := convert(*v)

	// Save current context for a potential call to NextPage
	c.params = params
	c.method = method
	c.totalPages = conv.Pages

	return convert(*v), nil
}

// NextPage retreives the next page in a paginated API call
func (c *PhotosClient) NextPage() (*PhotoList, error) {
	c.PaginationParams.Page = c.Page + 1
	if c.PaginationParams.Page > c.context.totalPages {
		return nil, ErrPaginatorExhausted
	}

	return c.Request(c.method, c.params)
}
