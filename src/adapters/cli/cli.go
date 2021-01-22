package cli

import (
	"errors"
	"strings"
	"unicode"

	"github.com/gitaction/docopt.go"
)

type command struct {
	usage string
	f     func(docopt.Opts)
}

var commands = make(map[string]*command)

func Register(name string, f func(opts docopt.Opts), usage string) error {
	if findCmd(name) {
		return errors.New("duplicate command")
	}
	commands[name] = &command{
		usage: strings.TrimLeftFunc(usage, unicode.IsSpace),
		f:     f,
	}
	return nil
}

func Run(name string, args []string) error {
	if !findCmd(name) {
		return errors.New("invalid command")
	}

	cmd := commands[name]
	parsedArgs, err := docopt.ParseArgs(cmd.usage, args, "0.0.1")
	if err != nil {
		return err
	}

	cmd.f(parsedArgs)
	return nil
}

func findCmd(name string) bool {
	_, found := commands[name]
	return found
}
