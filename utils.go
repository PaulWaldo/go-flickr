package flickr

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var ErrKeyNotInEnv = errors.New("no API Key in Environment")

func getApiKey(key, envFileName string) (string, error) {
	if key != "" {
		return key, nil
	}

	/* Check specified env file */
	if envFileName != "" {
		if _, envVarExists := os.LookupEnv(ApiKeyEnvVar); envVarExists {
			log.Printf("Value for env var %s already exists, but asking to read file %s.  File will not override", ApiKeyEnvVar, envFileName)
		}
		err := godotenv.Load(envFileName)
		fmt.Printf("Attempting to load env file %s\n", envFileName)
		if err != nil {
			return "", err
		}
	} else {
		/* Check for ./.env file */
		err := godotenv.Load()
		_, ok := err.(*os.PathError)
		if err != nil && !ok {
			return "", err
		}
	}
	if key, ok := os.LookupEnv(ApiKeyEnvVar); ok {
		fmt.Printf("Got key %s\n", key)
		return key, nil
	}
	fmt.Println("No key!!")
	return "", ErrKeyNotInEnv
}
