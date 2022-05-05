package cmdz

import (
	`os/exec`

	`github.com/photowey/hictl/pkg/logger`
)

func FormatSourceCode(filename string) {
	cmd := exec.Command("gofmt", "-w", filename)
	if err := cmd.Run(); err != nil {
		logger.Warnf("Error while running gofmt: %s", err)
	}
}

func ResolveDependencies() {
	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		logger.Warnf("running cmd: go mod tidy: %s", err)
	}
}
