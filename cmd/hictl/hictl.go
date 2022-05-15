package main

import (
	"log"

	"github.com/photowey/hictl"
	"github.com/photowey/hictl/cmd/internal/cmds/config"
	"github.com/photowey/hictl/pkg/logger"
	"github.com/spf13/cobra"

	"github.com/photowey/hictl/cmd/internal/cmds/schema"
	"github.com/photowey/hictl/cmd/internal/cmds/version"
	"github.com/photowey/hictl/cmd/internal/home"
)

func main() {
	logger.Infof("the hictl cmd pwd: %s", hictl.WorkDir)
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
		config.Cmd(),
		schema.Cmd(),
	)
	_ = cmd.Execute()
}
