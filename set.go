package main

import (
	"github.com/joho/godotenv"
	"github.com/simpleflags/cli/config"
	"log"
	"path"
)

type setCommand struct {
	ServerURL string `short:"s" long:"server" description:"Server URL address"`
	Account   string `short:"a" long:"acc" description:"Account identifier"`
	Project   string `short:"p" long:"project" description:"Project identifier"`
}

func (s setCommand) Execute(_ []string) error {
	sfDir, err := config.GetSimpleFlagsDir()
	if err != nil {
		log.Fatalf("Error getting data directory %v", err)
	}

	envs, err := godotenv.Read(path.Join(sfDir, ".env"))
	if err != nil {
		log.Println(err)
	}

	if s.ServerURL != "" {
		envs["SF_URL"] = s.ServerURL
	}

	if s.Account != "" {
		envs["SF_ACCOUNT"] = s.Account
	}

	if s.Project != "" {
		envs["SF_PROJECT"] = s.Project
	}

	return godotenv.Write(envs, path.Join(sfDir, ".env"))
}

func init() {
	sc := setCommand{}
	_, err := parser.AddCommand(
		"set",
		"Set environment variables",
		"Set environment variables like account, project (when value is empty then it will be removed)",
		&sc,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
