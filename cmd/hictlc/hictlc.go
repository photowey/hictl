package main

import (
	`fmt`
	`log`

	`github.com/hictl/cmd/internal/executor`
	`github.com/spf13/cobra`
)

func main() {
	log.SetFlags(0)
	cmd := &cobra.Command{
		Use: "hictlc",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("---- hictlc ----")
		},
	} // Replace entc command with hictlc
	cmd.AddCommand(
		executor.InitCmd(),
	)
	_ = cmd.Execute()
}
