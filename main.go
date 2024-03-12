package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"mydocker/cmd"
	"mydocker/config"
	"os"
)

/*
实现一个简单的容器运行时
参数：
 1. run 命令
 2. images 命令
 3. version 命令
*/
func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Printf("Failed to init config: %s\n", err.Error())
		os.Exit(1)
	}
	var rootCmd = &cobra.Command{
		Use:   "mydocker",
		Short: "mydocker is a simple container runtime implementation",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			_ = config.Init()
		},
	}
	rootCmd.AddCommand(
		cmd.Version,
		cmd.RunCmd,
		cmd.InitCmd,
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed to run app with %v: %s\n", os.Args, err.Error())
	}
}
