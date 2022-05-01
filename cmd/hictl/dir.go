package main

import (
	`bytes`
	`fmt`
	`os`
	`path/filepath`
	`strings`

	`github.com/hictl/cmd/internal/common/filez`
)

const (
	hictlHomeDir       = ".hictl"
	hictlConfig        = "hictl.json"
	hictlConfigContent = `{
  "databases": {}
}`
)

func checkHictlHome() {
	home, _ := os.UserHomeDir()
	hictlHome := filepath.Join(home, hictlHomeDir)
	if ok := filez.DirExists(hictlHome); !ok {
		if err := os.MkdirAll(hictlHome, os.ModePerm); err != nil {
			panic(fmt.Sprintf("mkdir hictl home dir:%s error:%v", hictlHome, err))
		}
	}

	if filez.FileNotExists(hictlHome, hictlConfig) {
		buf := bytes.NewBufferString(hictlConfigContent)
		hictlConfigFile := filepath.Join(hictlHome, strings.ToLower(hictlConfig))
		if err := os.WriteFile(hictlConfigFile, buf.Bytes(), 0644); err != nil {
			panic(fmt.Sprintf("writing file %s: %v", hictlConfigFile, err))
		}
	}

}