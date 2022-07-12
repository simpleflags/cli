package main

import (
	"fmt"
	"github.com/simpleflags/cli/config"
	"log"
	"os"
	"path"
)

type logoutCommand struct {
}

func (l logoutCommand) Execute(args []string) error {
	sfDir, err := config.GetSimpleFlagsDir()
	if err != nil {
		fmt.Printf("error getting simpleflags folder, err: %v", err)
	}
	err = os.Remove(path.Join(sfDir, "auth.data"))
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("logout error: %v\n", err)
		return err
	}
	fmt.Printf("Successfully logged out. Good Bye\n")
	return nil
}

func init() {
	lc := logoutCommand{}
	_, err := parser.AddCommand(
		"logout",
		"logout from the system",
		"Logout",
		&lc,
	)
	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
