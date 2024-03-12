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
		InitContainer()
	},
}
var (
	containerName string
	imageName     string
	interactive   bool
	tty           bool
	detach        bool
	port          int
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
	ExecContainer()
}

func ExecContainer() {
	cmd := os.Args[3]
	fmt.Println("will exec cmd=", cmd)
	err := syscall.Exec(cmd, os.Args[3:], os.Environ())
	if err != nil {
		fmt.Println("exec proc fail ", err)
		return
	}
	fmt.Println("forever exec it ")
	return
}
