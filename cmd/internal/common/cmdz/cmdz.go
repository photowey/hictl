package cmdz

import (
	`bufio`
	`bytes`
	`io`
	`log`
	`os`
	`os/exec`
)

func FormatCode(goFile string) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("gofmt", "-w", goFile)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Wait() failed with '%s'\n", err)
	}
}

func DeleteBlankLike(src string, dst string) {
	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0666)
	defer func(srcFile *os.File) {
		_ = srcFile.Close()
	}(srcFile)

	if err != nil {
		return
	}
	srcReader := bufio.NewReader(srcFile)
	destFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0666)
	defer func(destFile *os.File) {
		_ = destFile.Close()
	}(destFile)

	if err != nil {
		return
	}
	for {
		str, _ := srcReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		}
		if 0 == len(str) || str == "\r\n" {
			continue
		}
		_, _ = destFile.WriteString(str)
	}

	return
}
