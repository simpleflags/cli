package config

import (
	"os"
	"os/user"
	"path"
)

func GetSimpleFlagsDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	sfDir := path.Join(currentUser.HomeDir, ".simpleflags")

	if _, err = os.Stat(sfDir); os.IsNotExist(err) {
		err = os.Mkdir(sfDir, 0750)
		if err != nil {
			return "", err
		}
	}

	return sfDir, nil
}
