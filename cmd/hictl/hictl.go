package main

import (
	`log`

	`github.com/hictl/cmd/internal/executor`
	`github.com/spf13/cobra`
)

func main() {
	checkHictlHome()
	log.SetFlags(0)
	cmd := &cobra.Command{
		Use: "hictl",
		Run: func(cmd *cobra.Command, args []string) {
			// do noting
		},
	} // Replace ent command with hictl
	cmd.AddCommand(
		executor.InitCmd(),
	)
	_ = cmd.Execute()
}
