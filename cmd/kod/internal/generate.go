package internal

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
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
			if watch, _ := cmd.Flags().GetBool("watch"); watch {
				// Create new watcher.
				w := lo.Must(fsnotify.NewWatcher())
				defer w.Close()

				Watch(&watcher{w: w}, ".",
					func() { doGenerate(cmd, ".", args) },
					lo.Must(cmd.Flags().GetBool("verbose")),
				)
			}

			doGenerate(cmd, ".", args)
		}
	},
}

func doGenerate(cmd *cobra.Command, _ string, args []string) {
	startTime := time.Now()

	if s2i, _ := cmd.Flags().GetBool("struct2interface"); s2i {
		if err := Struct2Interface("."); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[struct2interface] %s \n", time.Since(startTime).String())
	}

	if err := Generate(".", args, Options{}); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("[generate] %s \n", time.Since(startTime).String())
}

func init() {
	generate.Flags().BoolP("struct2interface", "s", false, "generate interface from struct.")
	generate.Flags().BoolP("verbose", "v", false, "verbose mode.")
	generate.Flags().BoolP("watch", "w", false, "watch the changes of the files and regenerate the codes.")
	rootCmd.AddCommand(generate)
}
