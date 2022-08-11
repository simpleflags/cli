package main

import (
	"github.com/hashicorp/go-getter"
	"log"
)

type quickStartCommand struct {
	APIKey string `short:"k" long:"key" description:"Provide api key" env:"SF_API_KEY"`
	Lang   string `short:"l" long:"lang" description:"lang name (go, golang, js, javascript" choice:"go" choice:"golang" choice:"js" choice:"javascript" required:"true"`
	Dir    string `short:"d" long:"dir" description:"Output directory" required:"true"`
}

func (c quickStartCommand) Execute(args []string) error {
	return getter.Get("example", "git::https://github.com/simpleflags/quickstart.git//"+c.Lang)
}

func init() {
	sc := quickStartCommand{}
	_, err := parser.AddCommand(
		"quick",
		"quick start template project",
		"create quickstart project",
		&sc,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
