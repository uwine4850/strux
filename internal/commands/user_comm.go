package commands

import (
	"bufio"
	"fmt"
	"github.com/uwine4850/strux_api/services/protofiles/baseproto"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"strux/internal/apiutils"
	"syscall"
)

type UserCommands struct {
	User             string `short:"usr" long:"user" block:"1"`
	Register         string `short:"reg" long:"register"`
	UpdatePassword   string `short:"upd-pass" long:"update-password"`
	Delete           string `short:"del" long:"delete"`
	isRegister       bool
	isUpdatePassword bool
	isDelete         bool
}

func (uc *UserCommands) ExecUser() []string {
	return []string{uc.Register, uc.UpdatePassword, uc.Delete}
}

func (uc *UserCommands) ExecDelete() []string {
	uc.isDelete = true
	return []string{}
}

func (uc *UserCommands) ExecUpdatePassword() []string {
	uc.isUpdatePassword = true
	return []string{}
}

func (uc *UserCommands) ExecRegister() []string {
	uc.isRegister = true
	return []string{}
}

func (uc *UserCommands) OnFinish() {
	if uc.isRegister {
		response, err := register()
		if !response.Success {
			fmt.Println(response.Message)
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(response.Message)
		return
	}
	if uc.isUpdatePassword {
		response, err := updatePassword()
		if !response.Success {
			fmt.Println(response.Message)
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(response.Message)
		return
	}
	if uc.isDelete {
		response, err := deleteUser()
		if !response.Success {
			fmt.Println(response.Message)
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(response.Message)
		return
	}
}

func deleteUser() (*baseproto.BaseResponse, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("⚠️ Warning! Deleting a user will result in the deletion of all packages belonging to that user.")
	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Print("Enter password(hidden): ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Print("\n")
	apiForm := &apiutils.NewApiForm{
		Method: "DELETE",
		Url:    "http://localhost:3000/user-delete/",
		TextValues: map[string]string{
			"username": strings.TrimSpace(string(username)),
			"password": strings.TrimSpace(string(bytePassword)),
		},
		FileValues: nil,
	}
	baseResponse, _, err := apiForm.SendForm()
	if err != nil {
		return nil, err
	}
	return baseResponse, err
}

// updatePassword Updating the user's password.
// You must enter the username correctly to update it. Then enter the new password twice.
func updatePassword() (*baseproto.BaseResponse, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Print("Enter password(hidden): ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Print("\n")
	fmt.Print("Enter new password: ")
	newPassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Print("\n")
	fmt.Print("Enter new password again: ")
	newPasswordAgain, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Print("\n")
	// compare passwords
	if strings.TrimSpace(string(newPassword)) != strings.TrimSpace(string(newPasswordAgain)) {
		return &baseproto.BaseResponse{
			Message: "Passwords don`t match.",
			Success: false,
			Status:  baseproto.ResponseStatus_StatusWarning,
		}, nil
	}
	apiForm := &apiutils.NewApiForm{
		Method: "PUT",
		Url:    "http://localhost:3000/user-password-update/",
		TextValues: map[string]string{
			"username":    strings.TrimSpace(string(username)),
			"password":    strings.TrimSpace(string(bytePassword)),
			"newPassword": strings.TrimSpace(string(newPassword)),
		},
		FileValues: nil,
	}
	baseResponse, _, err := apiForm.SendForm()
	if err != nil {
		return nil, err
	}
	return baseResponse, err
}

// register Connecting to the api service. Sending a form with data about the new user.
func register() (*baseproto.BaseResponse, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Print("Enter password(hidden): ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Print("\n")
	fmt.Print("Enter password(again): ")
	bytePasswordAgain, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Print("\n")

	// compare passwords
	if string(bytePassword) != string(bytePasswordAgain) {
		return &baseproto.BaseResponse{
			Message: "Passwords don`t match.",
			Success: false,
			Status:  baseproto.ResponseStatus_StatusWarning,
		}, nil
	}
	apiForm := &apiutils.NewApiForm{
		Method:     "POST",
		Url:        "http://localhost:3000/create-user/",
		TextValues: map[string]string{"username": strings.TrimSpace(username), "password": strings.TrimSpace(string(bytePassword))},
		FileValues: nil,
	}
	baseResponse, _, err := apiForm.SendForm()
	if err != nil {
		return nil, err
	}
	return baseResponse, err
}
