package main

import (
	"fmt"
	"os"

	"github.com/dtan4/s3url/cli"
)

func main() {
	c := cli.New(os.Stdout, os.Stderr, Version, Revision)

	if err := c.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
