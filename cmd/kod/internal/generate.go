package internal

import (
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var generate = &cobra.Command{
	Use:   "generate",
	Short: "generate kod related codes for your kod application.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		{
			startTime := time.Now()

			lo.Must0(Generate(".", args, Options{}))

			fmt.Printf("[generate] %s \n", time.Since(startTime).String())
		}
	},
}

func init() {
	rootCmd.AddCommand(generate)
}
