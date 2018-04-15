package main

import (
	"os"

	"github.com/dtan4/s3url/cli"
)

func main() {
	c := cli.New(os.Stdout, os.Stderr, Version, Revision)

	os.Exit(c.Run(os.Args))
}
