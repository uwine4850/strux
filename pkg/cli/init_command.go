package cli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// CommandCollector contains all commands of a particular command structure
type CommandCollector struct {
	Commands            map[string][]string
	BlockCommand        string
	ParentStructCommand CommandInterface
}

func (cc CommandCollector) GetCommands() map[string][]string {
	return cc.Commands
}

func (cc CommandCollector) GetBlockCommand() string {
	return cc.BlockCommand
}

// InitCommand Initializes individual parts of the command structure.
// Initializes and returns a CommandCollector structure.
type InitCommand struct {
	Command CommandInterface
}

// Init Initializes a CommandCollector structure. Collection point for all helper methods.
func (mc *InitCommand) Init() CommandCollector {
	mc.fieldInit()

	rt := reflect.ValueOf(mc.Command).Type()
	cc := CommandCollector{}

	// processing each CommandInterface field
	for i := 0; i < rt.Elem().NumField(); i++ {
		field := rt.Elem().Field(i)
		// if tag exist
		if !reflect.DeepEqual(field.Tag, reflect.StructTag("")) {
			mc.collect(&field, &cc)
		}
	}
	//cc.SetParentStructCommand(mc.Command)
	if cc.BlockCommand != "" {
		d := &cc
		d.ParentStructCommand = mc.Command
		return cc
	} else {
		err := &ErrMissingBlockCommand{CommandStructName: reflect.TypeOf(mc.Command).Elem().Name()}
		panic(err)
	}
}

// fieldInit fills structure fields with their own names. Field = The name of this field.
// The method does not interact with the parent structure.
func (mc *InitCommand) fieldInit() {
	rt := reflect.ValueOf(mc.Command).Elem()

	// processing each CommandInterface field
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Type().Field(i)
		if reflect.DeepEqual(field.Type.Kind(), reflect.String) {
			if field.Tag != "" {
				rt.FieldByName(field.Name).SetString(field.Name)
			}
		}
	}
}

// collect Collects the fields of the command structure, reads the tags and packs them into a map.
func (mc *InitCommand) collect(field *reflect.StructField, cc *CommandCollector) {
	if len(cc.Commands) == 0 {
		cc.Commands = make(map[string][]string)
	}

	// example {"FName": ["-fn", "--fname"]}
	cc.Commands[field.Name] = append(cc.Commands[field.Name], mc.getShortCommand(field))
	cc.Commands[field.Name] = append(cc.Commands[field.Name], mc.getLongCommand(field))

	// set block command
	if _, ok := field.Tag.Lookup("block"); ok {
		cc.BlockCommand = mc.getBlockCommand(field)
	}
}

// getShortCommand Returns the tag whose name is equal to "short".
func (mc *InitCommand) getShortCommand(field *reflect.StructField) string {
	if _, ok := field.Tag.Lookup("short"); !ok {
		err := &ErrShortTagNotExist{CommandName: field.Name}
		panic(err)
	}
	return fmt.Sprint("-", strings.Trim(field.Tag.Get("short"), " "))
}

// getLongCommand Returns the tag whose name is equal to "short".
func (mc *InitCommand) getLongCommand(field *reflect.StructField) string {
	if _, ok := field.Tag.Lookup("long"); !ok {
		err := &ErrLongTagNotExist{CommandName: field.Name}
		panic(err)
	}
	return fmt.Sprint("--", strings.Trim(field.Tag.Get("long"), " "))
}

// getLongCommand Returns the tag whose name is equal to "block".
func (mc *InitCommand) getBlockCommand(field *reflect.StructField) string {
	block := field.Tag.Get("block")
	if val, err := strconv.Atoi(block); err != nil {
		panic(err)
	} else if val == 1 {
		return strings.Trim(field.Name, " ")
	}
	err := &ErrGettingBlockCommand{CommandName: field.Name}
	panic(err)
}
