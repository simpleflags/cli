package main

import (
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"io/ioutil"
	"log"
	"net/http"
)

type helpCommand struct {
	APIKey string `short:"k" long:"key" description:"Provide api key"`
	Lang   string `short:"l" long:"lang" description:"lang name (go, golang, js, javascript" choice:"go" choice:"golang" choice:"js" choice:"javascript"`
}

func (c helpCommand) Execute(args []string) error {
	var (
		err      error
		response *http.Response
	)
	switch c.Lang {
	case "go", "golang":
		response, err = http.Get("https://raw.githubusercontent.com/simpleflags/golang-server-sdk/main/README.md")
		if err != nil {
			return err
		}
	case "js", "javascript":
		response, err = http.Get("https://raw.githubusercontent.com/simpleflags/golang-server-sdk/main/README.md")
		if err != nil {
			return err
		}
	}
	if response.StatusCode == http.StatusOK {
		all, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		render := markdown.Render(string(all), 80, 6)
		fmt.Println(string(render))
	}
	return nil
}

func init() {
	hc := helpCommand{}
	_, err := parser.AddCommand(
		"help",
		"Show SDK integration instructions",
		"Show SDK integration instructions",
		&hc,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
