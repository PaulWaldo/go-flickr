package flickr

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var ErrKeyNotInEnv = errors.New("no API Key in Environment")

func GetApiKey(key, envFileName string) (string, error) {
	if key != "" {
		return key, nil
	}
	/* Check specified env file */
	if envFileName != "" {
		err := godotenv.Load(envFileName)
		if err != nil {
			return "", err
		}
	} else {
		/* Check for ./.env file */
		err := godotenv.Load()
		if err != nil {
			return "", err
		}
	}
	if key, ok := os.LookupEnv(ApiKeyEnvVar); ok {
		return key, nil
	}
	return "", ErrKeyNotInEnv
}
