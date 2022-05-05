package config

import (
	`io/ioutil`
	`path/filepath`
	`strings`

	`github.com/photowey/hictl`
	`github.com/photowey/hictl/cmd/internal/common/filez`
	`github.com/photowey/hictl/cmd/internal/common/helper`
	`github.com/photowey/hictl/pkg/color`
	`github.com/photowey/hictl/pkg/logger`
	`github.com/spf13/cobra`
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "show the hictl config info",
		Example: helper.Examples(
			"hictl config Example",
			"hictl config",
		),
		Run: func(cmd *cobra.Command, args []string) {
			hictlHome := hictl.HomeDir
			hictlConfigFile := filepath.Join(hictlHome, strings.ToLower(hictl.Config))
			if filez.FileExists(hictlHome, hictl.Config) {
				conf, err := ioutil.ReadFile(hictlConfigFile)
				if err != nil {
					logger.Errorf("parse the hictl config error:%s", err.Error())
					return
				}
				logger.Printf(color.GrayBold(string(conf)))
			} else {
				logger.Warnf("the hictl config file not exists")
			}
		},
	}

	return cmd
}
