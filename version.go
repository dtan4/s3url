package main

import (
	"fmt"
)

var (
	Version  string
	Revision string
)

func printVersion() {
	fmt.Println("s3url version " + Version + ", build " + Revision)
}
