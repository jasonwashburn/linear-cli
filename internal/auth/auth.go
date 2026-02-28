package auth

import (
	"errors"

	"github.com/spf13/viper"
)

// GetAPIKey returns the Linear API key from the environment.
// It returns an error if the key is not set.
func GetAPIKey() (string, error) {
	key := viper.GetString("LINEAR_API_KEY")
	if key == "" {
		return "", errors.New("LINEAR_API_KEY is not set; export it before running this command")
	}
	return key, nil
}
