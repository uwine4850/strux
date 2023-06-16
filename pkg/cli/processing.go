package cli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ExecCommand executes commands of the selected command structure
type ExecCommand struct {
	Command           CommandInterface
	CommandCollection CommandCollector
	cliCommands       []AcceptedCommand
	validator         ExecValidator
}

// SetCliCommand setting commands from the console with to this structure using AcceptedCommand
func (ec *ExecCommand) SetCliCommand(cliCommands []AcceptedCommand) *ExecCommand {
	ec.cliCommands = cliCommands
	return ec
}

// Start all commands
func (ec *ExecCommand) Start() {
	if err := ec.validator.CheckCommandExist(&ec.CommandCollection, &ec.cliCommands); err != nil {
		panic(err)
	}
	cNames := []string{
		ec.getBlockCommand(),
	}
	ec.Exec(cNames)

	// on finish commands chain
	ec.Command.OnFinish()
}

// Exec Executes commands from the console in the selected command structure.
// Commands are executed strictly in order.
// When the method finds a matching command, the exec method of the command structure is called,
// the name of which consists of the Exec prefix and the name of the field (command).
func (ec *ExecCommand) Exec(cNames []string) {
	var subCommands []string
	if len(ec.cliCommands) > 0 {
		acceptedCommand := ec.cliCommands[0]
		for cName := range cNames {
			flags := ec.CommandCollection.Commands[cNames[cName]]
			for f := range flags {
				// if struct_command item == console command
				if flags[f] == acceptedCommand.CommandName {
					callInput := make([]reflect.Value, 0)
					ec.parseAcceptedCommandValues(&acceptedCommand, &callInput)
					commandExecMethod := reflect.ValueOf(ec.Command).MethodByName(fmt.Sprint("Exec", cNames[cName]))

					// if command method not exist
					if !commandExecMethod.IsValid() {
						err := &ErrCommandMethodNotExist{MethodName: fmt.Sprint("Exec", cNames[cName])}
						panic(err)
					}

					if err := ec.validator.CheckLenMethod(&commandExecMethod, &acceptedCommand); err != nil {
						panic(err)
					}

					val := commandExecMethod.Call(callInput)[0].Interface().([]string)
					if reflect.DeepEqual(val, []string{}) {
						break
					} else {
						ec.removeCliCommandByName(acceptedCommand.CommandName)
					}
					ec.Exec(val)
				} else {
					// Commands that fail validation for one of two reasons:
					// 1. The command does not match the flag (-, --) but it is visible.
					// 2. The command is not valid.
					// In this block of code, 2 cases are tracked.
					notInSubCommands := true
					for i := 0; i < len(subCommands); i++ {
						if subCommands[i] == acceptedCommand.CommandName {
							notInSubCommands = false
						}
					}
					if notInSubCommands {
						c := ec.checkSubCommandExist(cNames, acceptedCommand.CommandName)
						if !c {
							err := &ErrSubcommandNotInThisContext{CommandName: acceptedCommand.CommandName}
							panic(err)
						} else {
							subCommands = append(subCommands, acceptedCommand.CommandName)
						}
					}
				}
			}
		}
	}
}

// checkSubCommandExist compares a list of commands with one single command.
// If the command is present (at least once) in the list, the method returns true.
func (ec *ExecCommand) checkSubCommandExist(cNames []string, consoleComm string) bool {
	for iName := 0; iName < len(cNames); iName++ {
		cc := ec.CommandCollection.Commands[cNames[iName]]
		for i := 0; i < len(cc); i++ {
			// One valid command equals the entire valid slice.
			if cc[i] == consoleComm {
				return true
			}
		}
	}
	return false
}

// removeCliCommandByName helper method to remove a console command from the list
func (ec *ExecCommand) removeCliCommandByName(name string) {
	var index []int
	for i := range ec.cliCommands {
		if name == ec.cliCommands[i].CommandName {
			index = append(index, i)
		}
	}
	for i := range index {
		ec.cliCommands = append(ec.cliCommands[:i], ec.cliCommands[i+1:]...)
	}
}

func (ec *ExecCommand) getBlockCommand() string {
	return ec.CommandCollection.BlockCommand
}

func (ec *ExecCommand) parseAcceptedCommandValues(aCommand *AcceptedCommand, callInput *[]reflect.Value) {
	for val := range aCommand.CommandValue {
		*callInput = append(*callInput, reflect.ValueOf(aCommand.CommandValue[val]))
	}
}

type AcceptedCommand struct {
	CommandName  string
	CommandValue []string
}

type ExecValidator struct {
}

func (ev ExecValidator) CheckLenMethod(method *reflect.Value, aCommand *AcceptedCommand) error {
	if method.Type().NumIn() != len(aCommand.CommandValue) {
		return &ErrMethodLenArguments{
			CommandName: aCommand.CommandName,
			CurrentLen:  strconv.Itoa(method.Type().NumIn()),
			ReceivedLen: strconv.Itoa(len(aCommand.CommandValue)),
		}
	}
	return nil
}

func (ev ExecValidator) CheckCommandExist(commandCollection *CommandCollector, aCommands *[]AcceptedCommand) error {
	var consoleCommands []string
	for j := 0; j < len(*aCommands); j++ {
		consoleCommands = append(consoleCommands, (*aCommands)[j].CommandName)
	}
	for _, commands := range commandCollection.Commands {
		for i := 0; i < len(commands); i++ {
			for j := 0; j < len(consoleCommands); j++ {
				if consoleCommands[j] == commands[i] {
					consoleCommands = append(consoleCommands[:j], consoleCommands[j+1:]...)
				}
			}
		}
	}
	if len(consoleCommands) > 0 {
		notExistCommands := strings.Join(consoleCommands, ", ")
		return &ErrCommandsNotExist{notExistCommands}
	}
	return nil
}
