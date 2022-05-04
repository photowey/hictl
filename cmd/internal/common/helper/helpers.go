package helper

import (
	`os/exec`
	`strings`

	`github.com/hictl/pkg/logger`
)

func MustCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func Examples(ex ...string) string {
	for i := range ex {
		ex[i] = "  " + ex[i]
	}
	return strings.Join(ex, "\n")
}

func GoVersion() string {
	var (
		cmdOut []byte
		err    error
	)

	if cmdOut, err = exec.Command("go", "version").Output(); err != nil {
		logger.Fatalf("There was an error running 'go version' command: %s", err)
	}
	return strings.Split(string(cmdOut), " ")[2]
}
