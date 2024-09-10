package internal

import (
	"context"
	"fmt"
	"path/filepath"
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
				startWatcher(cmd.Context(), cmd, args)
			}

			doGenerate(cmd, "./", args)
		}
	},
}

func startWatcher(ctx context.Context, cmd *cobra.Command, args []string) {
	// Create new watcher.
	w := lo.Must(fsnotify.NewWatcher())
	defer w.Close()

	Watch(&watcher{ctx: ctx, w: w}, ".",
		func(event fsnotify.Event) { doGenerate(cmd, filepath.Dir(event.Name), args) }, lo.Must(cmd.Flags().GetBool("verbose")),
	)
}

func doGenerate(cmd *cobra.Command, dir string, args []string) {
	startTime := time.Now()

	if s2i, _ := cmd.Flags().GetBool("struct2interface"); s2i {
		if err := Struct2Interface(cmd, dir); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[struct2interface] %s \n", time.Since(startTime).String())
	}

	startTime = time.Now()

	if err := Generate(dir, args, Options{}); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("[generate] %s \n", time.Since(startTime).String())
}

func init() {
	generate.Flags().BoolP("struct2interface", "s", true, "generate interface from struct.")
	generate.Flags().BoolP("verbose", "v", false, "verbose mode.")
	generate.Flags().BoolP("watch", "w", false, "watch the changes of the files and regenerate the codes.")
	rootCmd.AddCommand(generate)
}
