package filez

//
// filez
//

import (
	"os"
	"path/filepath"

	"github.com/photowey/hictl/cmd/internal/common/helper"
)

func exists(names ...string) bool {
	_, err := os.Stat(filepath.Join(names...))

	return err == nil || os.IsExist(err)
}

func DirExists(path string) bool {
	return exists(path)
}

func FileExists(target, name string) bool {
	return exists(target, name)
}

func FileNotExists(target, name string) bool {
	return !FileExists(target, name)
}

func Write(filename, content string) {
	f, err := os.Create(filename)
	helper.MustCheck(err)
	defer Close(f)
	_, err = f.WriteString(content)
	helper.MustCheck(err)
}

func Close(f *os.File) {
	err := f.Close()
	helper.MustCheck(err)
}
