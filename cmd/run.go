package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mydocker/ns"
	"mydocker/tools"
	"os"
	"os/exec"
	"syscall"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "run a container",
	Run: func(cmd *cobra.Command, args []string) {
		RunContainer()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if err := ns.DeleteMntNameSpace(containerName); err != nil {
			return
		}
	},
}
var (
	containerName string
	imageName     string
	interactive   bool
	tty           bool
	detach        bool
	port          string
)

func init() {
	RunCmd.Flags().StringVarP(&containerName, "name", "n", tools.GenerateDefaultName(), "container name")
	RunCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "interactive")
	RunCmd.Flags().BoolVarP(&tty, "tty", "t", false, "tty")
	RunCmd.Flags().BoolVarP(&detach, "detach", "d", false, "detach")
	RunCmd.MarkFlagRequired("name")
}
func RunContainer() {
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		fmt.Println("get init process error ", err)
		return
	}
	os.Args[1] = "init"
	cmd := exec.Command(initCmd, os.Args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("init proc end", initCmd)
	return
}
