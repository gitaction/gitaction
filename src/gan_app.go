package main

import (
	"fmt"
	"github.com/sunwei/gitaction/src/adapters/cli"
)

func main() {
	if err := cli.Process(); err != nil {
		fmt.Printf("%+v", err)
	}
}
