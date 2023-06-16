package cli

import "fmt"

type ErrCommandMethodNotExist struct {
	MethodName string
}

func (e *ErrCommandMethodNotExist) Error() string {
	return fmt.Sprintf("An error occurred while calling the %s method, it may not exist.", e.MethodName)
}

type ErrCommandsNotExist struct {
	CommandNames string
}

func (e *ErrCommandsNotExist) Error() string {
	return fmt.Sprintf("Commands %s not exist.", e.CommandNames)
}

type ErrMethodLenArguments struct {
	CommandName string
	CurrentLen  string
	ReceivedLen string
}

func (e *ErrMethodLenArguments) Error() string {
	return fmt.Sprintf("Command \"%s\" accept only %s arguments, received %s", e.CommandName, e.CurrentLen, e.ReceivedLen)
}

type ErrSubcommandNotInThisContext struct {
	CommandName string
}

func (e *ErrSubcommandNotInThisContext) Error() string {
	return fmt.Sprintf("The %s subcommand cannot be used in this context.", e.CommandName)
}

type ErrBlockCommandsDuplication struct {
	s string
}

func (e *ErrBlockCommandsDuplication) Error() string {
	return e.s
}

type ErrBlockCommandNotExist struct {
	CommandName string
}

func (e *ErrBlockCommandNotExist) Error() string {
	return fmt.Sprintf("Block command %s not exist.", e.CommandName)
}

type ErrMissingBlockCommand struct {
	CommandStructName string
}

func (e *ErrMissingBlockCommand) Error() string {
	return fmt.Sprintf("The tag block:\"1\" is missing from the %s.", e.CommandStructName)
}

type ErrShortTagNotExist struct {
	CommandName string
}

func (e *ErrShortTagNotExist) Error() string {
	return fmt.Sprintf("The \"%s\" command does not have a short tag.", e.CommandName)
}

type ErrLongTagNotExist struct {
	CommandName string
}

func (e *ErrLongTagNotExist) Error() string {
	return fmt.Sprintf("The \"%s\" command does not have a long tag.", e.CommandName)
}

type ErrGettingBlockCommand struct {
	CommandName string
}

func (e *ErrGettingBlockCommand) Error() string {
	return fmt.Sprintf("Error while getting command block for %s field.", e.CommandName)
}
