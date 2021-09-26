package flickr

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Params map[string]string

type Client struct {
	Key   string
	Token string
	Sig   string
}

const ApiKeyEnvVar = "FLICKR_API_KEY"

func NewClient(apiKey string, envFileName string) *Client {
	if apiKey != "" {
		return &Client{Key: apiKey}
	}
	var err error

	/* Check specified env file */
	if envFileName != "" {
		err = godotenv.Load(envFileName)
		if err != nil {
			log.Fatalf("Error loading env file %s", envFileName)
		}
	} else {
		/* Check for ./.env file */
		godotenv.Load()
	}

	if key, ok := os.LookupEnv(ApiKeyEnvVar); ok {
		return &Client{Key: key}
	}

	panic("Unable to get API Key")
}

func (client *Client) Request(method string, params Params) ([]byte, error) {
	url := fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.%s&api_key=%s&format=json&nojsoncallback=1", method, client.Key)
	if len(client.Token) > 0 {
		url = fmt.Sprintf("%s&auth_token=%s", url, client.Token)
	}

	if len(client.Sig) > 0 {
		url = fmt.Sprintf("%s&auth_sig=%s", url, client.Sig)
	}

	for key, value := range params {
		url = fmt.Sprintf("%s&%s=%s", url, key, value)
	}

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

type PaginatedClient struct {
	*Client
	Page       int `json:"page"`
	NumPages   int `json:"pages"`
	NumPerPage int `json:"perpage"`
	Total      int `json:"total"`
	CurPage    int
}

// NewPaginatedClient creates a Client that provides paginated results,
// numPerPage at a time
func NewPaginatedClient(apiKey string, envFileName string, numPerPage, page int) *PaginatedClient {
	p := &PaginatedClient{}
	p.Client = NewClient(apiKey, envFileName)
	p.NumPerPage = numPerPage
	p.Page = page
	return p
}

// NewDefaultPaginatedClient creates a PaginatedClient providing pages of 100 items starting at page 1
func NewDefaultPaginatedClient(apiKey string, envFileName string) *PaginatedClient {
	return NewPaginatedClient(apiKey, envFileName, 100, 1)
}

func (client *PaginatedClient) Request(method string, params Params) ([]byte, error) {
	params["per_page"] = strconv.Itoa(client.NumPerPage)
	params["page"] = strconv.Itoa(client.Page)
	return client.Client.Request(method, params)
}
