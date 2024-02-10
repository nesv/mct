package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/nesv/mct/microcode"
)

func main() {
	fmt.Println("Miek's Configuration Tool")

	cmds, err := microcode.ParseCommands(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', 0)
	defer tw.Flush()
	for _, cmd := range cmds {
		parts := []string{
			fmt.Sprintf("%8s", cmd.Name),
			strings.Join(cmd.Args, " "),
		}
		if !cmd.Action.IsZero() {
			parts = append(parts,
				"&&",
				cmd.Action.Name,
				strings.Join(cmd.Action.Args, " "),
			)
		}
		fmt.Fprintln(tw, strings.Join(parts, "\t"))
	}
}
