//go:build unit
// +build unit

package flickr

import (
	"testing"
)

func TestNewClientSpecifiedApiKey(t *testing.T) {
	expectedApiKey := "my-api-key"
	sut := &Client{Key: expectedApiKey}
	if sut.Key != expectedApiKey {
		t.Errorf("Expecting API Key '%s', but got '%s'", expectedApiKey, sut.Key)
	}
}

// func TestNewClientEnvFile(t *testing.T) {
// 	// Temporarily unset any API Key that is set
// 	apiKey, keyExists := os.LookupEnv(ApiKeyEnvVar)
// 	if keyExists {
// 		os.Unsetenv(ApiKeyEnvVar)
// 		t.Cleanup(func() {
// 			os.Setenv(ApiKeyEnvVar, apiKey)
// 		})
// 	}

// 	// Create a temporary env file
// 	file, err := ioutil.TempFile("", "go-test")
// 	if err != nil {
// 		t.Errorf("Unable to create temp env file: %s", err)
// 	}
// 	defer os.Remove(file.Name())

// 	// Write API key to file
// 	expectedApiKey := "my-api-key"
// 	file.WriteString(fmt.Sprintf("%s=%s", ApiKeyEnvVar, expectedApiKey))
// 	file.Close()

// 	sut, err := NewClientEnvFile(file.Name())
// 	if err != nil {
// 		t.Fatalf("Unable to create client: %s", err)
// 	}
// 	if sut.Key != expectedApiKey {
// 		t.Errorf("Expecting API Key '%s', but got '%s'", expectedApiKey, sut.Key)
// 	}
// }
