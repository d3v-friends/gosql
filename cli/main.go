package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var cmd = &cobra.Command{
		Use: "gosql",
	}

	cmd.AddCommand(CmdReader())
	cmd.AddCommand(CmdMigrate())

	var err = cmd.Execute()
	if err != nil {
		panic(err)
	}
}
