package flickr

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Params map[string]string

type Client struct {
	Key   string
	Token string
	Sig   string
}

// ApiKeyEnvVar is the name of the environment variable that is search for the Flickr API Key
const ApiKeyEnvVar = "FLICKR_API_KEY"

// NewClient creates a client that can access the Flickr API,
// attempting fo fetch the API Key from the file ./.env
func NewClient() (*Client, error) {
	return NewClientEnvFile("")
}

// NewClientEnvFile creates a client that can access the Flickr API,
// attempting fo fetch the API Key first from an environment file specified in envFileName
// then from the file ./.env
func NewClientEnvFile(envFileName string) (*Client, error) {
	/* Check specified env file */
	if envFileName != "" {
		err := godotenv.Load(envFileName)
		if err != nil {
			log.Fatalf("Error loading env file %s", envFileName)
		}
	} else {
		/* Check for ./.env file */
		godotenv.Load()
	}

	return NewClientEnvVar()
}

// NewClientEnvFile creates a client that can access the Flickr API,
// by searching for an environment variable named by ApiKeyEnvVar
func NewClientEnvVar() (*Client, error) {
	if key, ok := os.LookupEnv(ApiKeyEnvVar); ok {
		return &Client{Key: key}, nil
	}
	return nil, fmt.Errorf("API Key '%s' not found in environment", ApiKeyEnvVar)
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
	if response.StatusCode < http.StatusOK || response.StatusCode > http.StatusPermanentRedirect {
		return nil, fmt.Errorf("http status %s", response.Status)
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
