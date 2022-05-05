package hictl

import (
	`os`
	`os/user`
	`path/filepath`
)

const (
	Version string = "1.0.0"
	Home    string = ".hictl"
	Config         = "hictl.json"
)

var (
	Usr, _  = user.Current()
	HomeDir = filepath.Join(Usr.HomeDir, "/", Home)
	WorkDir = pwd()
	GoPath  = os.Getenv("GOPATH")
)

func pwd() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}

	return ""
}
