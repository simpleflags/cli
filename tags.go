package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

type tagsCommand struct {
	Account string `short:"a" long:"acc" description:"Account identifier" env:"SF_ACCOUNT"`
	Project string `short:"p" long:"project" description:"Project identifier" env:"SF_PROJECT"`
}

func (t tagsCommand) Execute(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	tags, err := api.GetTags(ctx, t.Account, t.Project, args...)
	if err != nil {
		return err
	}
	fmt.Println(strings.Join(tags, ", "))
	return nil
}

func init() {
	tc := tagsCommand{}
	_, err := parser.AddCommand(
		"tags",
		"List all tags in project",
		"List all tags in project",
		&tc,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
