package main

import (
	"fmt"
	"github.com/r3labs/sse/v2"
	"log"
	"os"
)

type streamCommand struct {
}

func (c streamCommand) Execute(_ []string) error {
	client := sse.NewClient("https://64a55c46.fanoutcdn.com/api/stream", func(c *sse.Client) {
		c.Headers["API-Key"] = os.Getenv("SF_API_KEY")
	})
	return client.Subscribe("", func(msg *sse.Event) {
		// Got some data!
		fmt.Println(string(msg.Data))
	})
}

func init() {
	_, err := parser.AddCommand(
		"stream",
		"Check stream connection",
		"Check stream connection",
		&streamCommand{},
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
