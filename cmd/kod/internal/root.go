/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package internal

import (
	"fmt"
	"runtime/debug"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kod",
	Short: "A powerful tool for writing kod applications.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if version, _ := cmd.Flags().GetBool("version"); version {
			info, ok := debug.ReadBuildInfo()
			if ok {
				fmt.Printf("%s %s\n", info.Main.Path, info.Main.Version)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	lo.Must0(rootCmd.Execute())
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kod.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Help message for toggle")
}
