package hictl

import (
	`os`
	`os/user`
	`path/filepath`
)

const (
	Version string = "1.0.0"
)
const (
	HictlHome = ".hictl"
)

var (
	Usr, _     = user.Current()
	HomeDir    = filepath.Join(Usr.HomeDir, "/", HictlHome)
	CurrentDir = pwd()
	GoPath     = os.Getenv("GOPATH")
)

func pwd() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}

	return ""
}
