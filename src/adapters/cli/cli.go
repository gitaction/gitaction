package cli

import (
	"errors"
	"github.com/gitaction/docopt.go"
)

var commands = make(map[string]func(docopt.Opts))

func Process() error {
	usage := `Usage:
  gan serve [--port=<port>]
  gan --help | --version`
	
	opts, err := docopt.ParseArgs(usage, nil, "0.0.1rc")
	if err != nil {
		return err
	}
	if err := Run(opts); err != nil {
		return err
	}
	return nil
}

func Register(name string, f func(opts docopt.Opts)) error {
	if findCmd(name) {
		return errors.New("duplicate command")
	}
	commands[name] = f
	return nil
}

func Run(opts docopt.Opts) error {
	commandNames := getCommandName()
	for _, name := range commandNames {
		if opts[name] == true {
			if !findCmd(name) {
				return errors.New("invalid command")
			}
			f := commands[name]
			f(opts)
			break
		}
	}
	return nil
}

func getCommandName() []string {
	var commandNames []string
	for name := range commands {
		commandNames = append(commandNames, name)
	}
	return commandNames
}

func findCmd(name string) bool {
	_, found := commands[name]
	return found
}
