package cli

import (
	"fmt"
	"github.com/gitaction/docopt.go"
)

func init() {
	err := Register("serve", runServe)
	if err != nil {
		fmt.Printf("Register command 'serve' error failed with error : %+v\n", err)
	}
}

func runServe(opts docopt.Opts) {
	fmt.Printf("%+v\n", opts["--port"])
	fmt.Println("hello world")
}
