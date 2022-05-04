package cmdz

import (
	`os/exec`

	`github.com/hictl/pkg/logger`
)

func FormatSourceCode(filename string) {
	cmd := exec.Command("gofmt", "-w", filename)
	if err := cmd.Run(); err != nil {
		logger.Warnf("Error while running gofmt: %s", err)
	}
}
