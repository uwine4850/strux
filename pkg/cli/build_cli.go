package cli

import "fmt"

type Build struct {
	Commands          *[]CommandInterface
	ConsoleArgs       []AcceptedCommand
	commandCollection []CommandCollector
}

// ExecBuild with the ExecCommand structure, runs ONE selected command structure
func (b *Build) ExecBuild() {
	b.initCommands()
	currentCommandCollection := b.getCurrentCollection()
	execCommand := ExecCommand{
		Command:           currentCommandCollection.ParentStructCommand,
		CommandCollection: currentCommandCollection,
	}
	execCommand.SetCliCommand(b.ConsoleArgs)
	execCommand.Start()
}

// initCommands using the structure InitCommand initializes all command structures
func (b *Build) initCommands() {
	for command := range *b.Commands {
		e := InitCommand{Command: (*b.Commands)[command]}
		b.commandCollection = append(b.commandCollection, e.Init())
	}
}

// getCurrentCollection finds a matching structure using the blocking command(block:"1") in the
// Command Collector's structure slice and returns it.
func (b *Build) getCurrentCollection() CommandCollector {
	block := make([]string, 0)
	for i := range b.commandCollection {
		block = append(block, b.commandCollection[i].BlockCommand)
	}
	_, err := b.checkDuplicateBlockCommands(block)
	if err != nil {
		panic(err)
	}
	for collector := range b.commandCollection {
		for i := 0; i < len(block); i++ {
			structCommandBlockList := b.commandCollection[collector].Commands[block[i]]
			if len(structCommandBlockList) != 0 {
				for iC := range structCommandBlockList {
					if structCommandBlockList[iC] == b.ConsoleArgs[0].CommandName {
						return b.commandCollection[collector]
					}
				}
			}
		}
	}
	err = &ErrBlockCommandNotExist{CommandName: b.ConsoleArgs[0].CommandName}
	panic(err.Error())
}

// checkDuplicateBlockCommands checks if there are duplicate blocking commands in the current command structures.
func (b *Build) checkDuplicateBlockCommands(blockFieldsName []string) (string, error) {
	var allBlockList []string
	// get block commands
	for collector := range b.commandCollection {
		for i := 0; i < len(blockFieldsName); i++ {
			structCommandBlockList := b.commandCollection[collector].Commands[blockFieldsName[i]]
			for j := 0; j < len(structCommandBlockList); j++ {
				allBlockList = append(allBlockList, structCommandBlockList[j])
			}
		}
	}
	// check for duplicate block command
	for i := 0; i < len(allBlockList); i++ {
		newAllBlockList := allBlockList[i+1:]
		for j := range newAllBlockList {
			if allBlockList[i] == newAllBlockList[j] {
				mErr := ErrBlockCommandsDuplication{s: fmt.Sprintf("Block command %s duplicated.", allBlockList[i])}
				return allBlockList[i], &mErr
			}
		}
	}
	return "", nil
}
