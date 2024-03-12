package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mydocker/config"
)

var Version = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mydocker version:  " + config.Cfg.Version)
	},
}
