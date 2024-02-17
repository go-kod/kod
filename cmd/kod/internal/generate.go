package internal

import (
	"context"
	"fmt"
	"time"

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
			ctx := cmd.Context()
			if timeout, _ := cmd.Flags().GetDuration("timeout"); timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}

			if watch, _ := cmd.Flags().GetBool("watch"); watch {
				Watch(ctx, ".", func() {
					doGenerate(cmd, ".", args)
				})
			}

			doGenerate(cmd, ".", args)
		}
	},
}

func doGenerate(cmd *cobra.Command, dir string, args []string) {
	startTime := time.Now()

	if s2i := cmd.Flag("struct2interface").Changed; s2i {
		Struct2Interface(cmd, ".")
		fmt.Printf("[struct2interface] %s \n", time.Since(startTime).String())
	}

	err := Generate(".", args, Options{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("[generate] %s \n", time.Since(startTime).String())
}

func init() {
	generate.Flags().DurationP("timeout", "t", 0, "timeout for the generation.")
	generate.Flags().BoolP("struct2interface", "s", false, "generate interface from struct.")
	generate.Flags().BoolP("watch", "w", false, "watch the changes of the files and regenerate the codes.")
	rootCmd.AddCommand(generate)
}
