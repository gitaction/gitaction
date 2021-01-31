package main

import (
	"fmt"
	"github.com/gitaction/docopt.go"
)

func main() {
	usage := `usage: gan [-h|--help] [--version] <command> [<args>...]

Options:
  -h, --help                 Show this message
  --version                  Show current version

Commands:
  help                       Show usage for a specific command
  serve                      Listen and serve http service

See 'gan <command>' for more information on a specific command.
`
	parsedArgs, _ := docopt.ParseArgs(usage, nil, "0.0.1")
	cmd, _ := parsedArgs.String("<command>")
	
	if cmd == "help" {
		fmt.Print(usage)
		return
	}
	
}
