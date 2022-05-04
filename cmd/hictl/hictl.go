package main

import (
	`log`

	`github.com/hictl/cmd/internal/cmds/schema`
	`github.com/hictl/cmd/internal/cmds/version`
	`github.com/hictl/cmd/internal/home`
	`github.com/spf13/cobra`
)

func main() {
	home.HictlHome()
	log.SetFlags(0)
	cmd := &cobra.Command{
		Use: "hictl",
		Run: func(cmd *cobra.Command, args []string) {
			// do noting
		},
	} // Replace ent command with hictl
	cmd.AddCommand(
		version.Cmd(),
		schema.Cmd(),
	)
	_ = cmd.Execute()
}
