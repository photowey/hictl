package home

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/photowey/hictl"
	"github.com/photowey/hictl/cmd/internal/common/filez"
)

const (
	hictlConfig        = "hictl.json"
	hictlConfigContent = `{
  "databases": {}
}`
)

func HictlHome() {
	hictlHome := hictl.HomeDir
	if ok := filez.DirExists(hictlHome); !ok {
		if err := os.MkdirAll(hictlHome, os.ModePerm); err != nil {
			panic(fmt.Sprintf("mkdir hictl home dir:%s error:%v", hictlHome, err))
		}
	}

	if filez.FileNotExists(hictlHome, hictlConfig) {
		buf := bytes.NewBufferString(hictlConfigContent)
		hictlConfigFile := filepath.Join(hictlHome, strings.ToLower(hictlConfig))
		if err := os.WriteFile(hictlConfigFile, buf.Bytes(), 0o644); err != nil {
			panic(fmt.Sprintf("writing file %s: %v", hictlConfigFile, err))
		}
	}
}
