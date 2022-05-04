package version

import (
	`bytes`
	`os`

	`github.com/hictl/cmd/internal/cmds/version/banner`
	`github.com/hictl/cmd/internal/common/helper`
	`github.com/hictl/pkg/color`
	`github.com/spf13/cobra`
)

const (
	BannerTemplate string = `  _     _      _   _ 
 | |   (_)    | | | |
 | |__  _  ___| |_| |
 | '_ \| |/ __| __| |
 | | | | | (__| |_| |
 |_| |_|_|\___|\__|_| 
GoVersion: {{ .GoVersion }}
GOOS: {{ .GOOS }}
GOARCH: {{ .GOARCH }}
NumCPU: {{ .NumCPU }}
GOPATH: {{ .GOPATH }}
GOROOT: {{ .GOROOT }}
Compiler: {{ .Compiler }}
Version: v{{ .Version }}
`
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show the current hictl version",
		Example: helper.Examples(
			"hictl version Example",
			"hictl version",
		),
		Run: func(cmd *cobra.Command, args []string) {
			output := color.NewColorWriter(os.Stdout)
			banner.Init(output, bytes.NewBufferString(color.YellowBold(BannerTemplate)))
		},
	}

	return cmd
}
