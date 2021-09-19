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
	}

	/* Check for ./.env file */
	godotenv.Load()

	if key, found := os.LookupEnv(ApiKeyEnvVar); found {
		return &Client{Key: key}
	}

	panic("Unable to get API Key")
}

func (client *Client) Request(method string, params Params) ([]byte, error) {
	url := fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.%s&api_key=%s&format=json&nojsoncallback=1", method, client.Key)
	for k, v := range params {
		url += fmt.Sprintf("&%s=%s", k, v)
	}

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
