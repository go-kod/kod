package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/dominikbraun/graph/draw"
	"github.com/go-kod/kod/internal/callgraph"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var callgraphCmd = &cobra.Command{
	Use:   "callgraph",
	Short: "generate kod callgraph for your kod application.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		if len(args) == 0 {
			cmd.PrintErr("please input the binary filepath")
			return
		}

		g := lo.Must(callgraph.ReadComponentGraph(args[0]))
		o := lo.Must(cmd.Flags().GetString("o"))

		file := lo.Must(os.Create(o))
		lo.Must0(draw.DOT(g, file))

		fmt.Printf("[callgraph] %s \n", time.Since(startTime).String())
	},
}

func init() {

	callgraphCmd.PersistentFlags().String("o", "my-graph.dot", "output file name")

	rootCmd.AddCommand(callgraphCmd)
}
