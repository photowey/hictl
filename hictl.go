package hictl

import (
	`os`
	`os/user`
	`path/filepath`
)

const (
	Home    string = ".hictl"
	Version string = "1.0.0"
)

var (
	Usr, _     = user.Current()
	HomeDir    = filepath.Join(Usr.HomeDir, "/", Home)
	CurrentDir = pwd()
	GoPath     = os.Getenv("GOPATH")
)

func pwd() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}

	return ""
}
