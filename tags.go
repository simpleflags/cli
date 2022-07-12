package main

import "log"

type tagsCommand struct {
}

func (t tagsCommand) Execute(_ []string) error {
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
