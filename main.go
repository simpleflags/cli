/*
Copyright © 2022 Enver Bisevac <enver@bisevac.com>

*/
package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"github.com/simpleflags/cli/config"
	"log"
	"os"
	"path"
)

var (
	parser = flags.NewParser(nil, flags.Default)
)

func main() {
	sfDir, err := config.GetSimpleFlagsDir()
	if err != nil {
		log.Fatalf("Error getting data directory %v", err)
	}

	err = godotenv.Load(path.Join(sfDir, ".env"))
	if err != nil {
		log.Printf("Error loading .env file")
	}

	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
		default:
			os.Exit(1)
		}
	}
}
