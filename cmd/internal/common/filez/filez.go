package filez

//
// filez
//

import (
	`os`
	`path/filepath`
)

func exists(names ...string) bool {
	var _, err = os.Stat(filepath.Join(names...))

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
