package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mydocker/ns"
	"mydocker/tools"
	"os"
	"syscall"
)

var InitCmd = &cobra.Command{
	Use:   "init [imageName]",
	Short: "init a container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imageName = args[0]
		command = args[1:]
		InitContainer()
		ExecContainer()
	},
}
var (
	containerName string
	imageName     string
	interactive   bool
	tty           bool
	detach        bool
	port          int
	command       []string
)

func init() {
	InitCmd.Flags().StringVarP(&containerName, "name", "n", tools.GenerateDefaultName(), "container name")
	InitCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "interactive")
	InitCmd.Flags().BoolVarP(&tty, "tty", "t", false, "tty")
	InitCmd.Flags().BoolVarP(&detach, "detach", "d", false, "detach")
	InitCmd.Flags().IntVarP(&port, "port", "p", 8000, "port")
}

func InitContainer() {
	if err := ns.InitMntNameSpace(imageName, containerName); err != nil {
		fmt.Println("create mnt namespace fail ", err)
		return
	}
	if err := ns.MountProc(); err != nil {
		fmt.Println("mount proc fail ", err)
		return
	}
}

func ExecContainer() {
	fmt.Println("will exec cmd=", command)
	err := syscall.Exec(command[0], command, os.Environ())
	if err != nil {
		fmt.Println("exec proc fail ", err)
		return
	}
	fmt.Println("forever exec it ")
	return
}
