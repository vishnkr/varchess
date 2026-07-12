package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vc-server/cmd/vc-cli/commands"
)

func main() {
	root := &cobra.Command{
		Use:   "vc-cli",
		Short: "VarChess CLI — local engine REPL and multiplayer client",
	}

	root.AddCommand(commands.EngineCmd())
	root.AddCommand(commands.PlayCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
