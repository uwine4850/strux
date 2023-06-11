package main

import (
	"os"
	"strux/internal/commands"
	"strux/pkg/cli"
)

func main() {
	parser := cli.Parser{
		CommandArgs: os.Args[1:],
	}
	build := cli.Build{
		Commands: &[]cli.CommandInterface{
			&commands.CreateCommand{},
			&commands.InitCommand{},
			&commands.InfoCommand{},
		},
		ConsoleArgs: parser.Parse(),
	}
	build.ExecBuild()
}
