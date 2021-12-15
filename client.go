package flickr

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Params map[string]string

type Client struct {
	Key   string
	Token string
	Sig   string
	URL   string
}

// ApiKeyEnvVar is the name of the environment variable that is search for the Flickr API Key
const ApiKeyEnvVar = "FLICKR_API_KEY"

const flickrURL = "https://api.flickr.com/services/rest"

// NewClient creates a client that can access the Flickr API,
// attempting fo fetch the API Key from the file ./.env
func NewClient() (*Client, error) {
	return NewClientEnvFile("")
}

// NewClientEnvFile creates a client that can access the Flickr API,
// attempting fo fetch the API Key first from an environment file specified in envFileName
// then from the file ./.env
func NewClientEnvFile(envFileName string) (*Client, error) {
	key, err := GetApiKey("", envFileName)
	if err != nil {
		return nil, err
	}
	return &Client{Key: key, URL: flickrURL}, nil
}

// NewClientEnvFile creates a client that can access the Flickr API,
// by searching for an environment variable named by ApiKeyEnvVar
func NewClientEnvVar() (*Client, error) {
	key, err := GetApiKey(ApiKeyEnvVar, "")
	if err != nil {
		return nil, err
	}
	return &Client{Key: key, URL: flickrURL}, nil
}

func (c *Client) Request(method string, params Params) ([]byte, error) {
	url := fmt.Sprintf("%s/?method=flickr.%s&api_key=%s&format=json&nojsoncallback=1", c.URL, method, c.Key)
	if len(c.Token) > 0 {
		url = fmt.Sprintf("%s&auth_token=%s", url, c.Token)
	}

	if len(c.Sig) > 0 {
		url = fmt.Sprintf("%s&auth_sig=%s", url, c.Sig)
	}

	for key, value := range params {
		url = fmt.Sprintf("%s&%s=%s", url, key, value)
	}

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode > http.StatusPermanentRedirect {
		return nil, fmt.Errorf("http status %s", response.Status)
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
