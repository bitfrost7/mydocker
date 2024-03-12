package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mydocker/ns"
	"os"
	"os/exec"
	"syscall"
)

var RunCmd = &cobra.Command{
	Use:                "run",
	Short:              "run a container",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		RunContainer()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		for i, arg := range args {
			if arg == "-n" || arg == "--name" {
				ClearContainer(args[i+1])
				break
			}
		}
	},
}

func RunContainer() {
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		fmt.Println("get init process error ", err)
		return
	}
	os.Args[1] = "init"
	fmt.Println("command args:", os.Args)
	cmd := exec.Command(initCmd, os.Args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if cmdErr := cmd.Run(); cmdErr != nil {
		fmt.Println(cmdErr)
	}
	fmt.Println("init proc end", initCmd)
	return
}

// ClearContainer 清理容器
func ClearContainer(containerName string) {
	fmt.Println("clear container")
	if err := ns.DeleteMntNameSpace(containerName); err != nil {
		fmt.Println("delete mnt namespace fail ", err)
		return
	}
}
