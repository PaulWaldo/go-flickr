package flickr

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func ExampleRequest() {
	client := NewClient("", "")

	resp, err := client.Request("people.findByUsername", Params{
		"username": "azerbike",
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v", string(resp))
	// Output: {"user":{"id":"98269877@N00","nsid":"98269877@N00","username":{"_content":"azerbike"}},"stat":"ok"}
}

func TestNewClientSpecifiedApiKey(t *testing.T) {
	expectedApiKey := "my-api-key"
	sut := NewClient(expectedApiKey, "")
	if sut.Key != expectedApiKey {
		t.Errorf("Expecting API Key '%s', but got '%s'", expectedApiKey, sut.Key)
	}
}

func TestNewClientEnvFile(t *testing.T) {
	// Create a temporary env file
	file, err := ioutil.TempFile("", "go-test")
	if err != nil {
		t.Errorf("Unable to create temp env file: %s", err)
	}
	defer os.Remove(file.Name())

	// Write API key to file
	expectedApiKey := "my-api-key"
	file.WriteString(fmt.Sprintf("%s=%s", ApiKeyEnvVar, expectedApiKey))
	file.Close()

	sut := NewClient("", file.Name())
	if sut.Key != expectedApiKey {
		t.Errorf("Expecting API Key '%s', but got '%s'", expectedApiKey, sut.Key)
	}
}
