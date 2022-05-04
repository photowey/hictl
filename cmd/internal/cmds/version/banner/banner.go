package banner

import (
	`io`
	`io/ioutil`
	`runtime`
	`text/template`

	`github.com/hictl`
	`github.com/hictl/cmd/internal/common/helper`
	`github.com/hictl/pkg/color`
	`github.com/hictl/pkg/logger`
)

type RuntimeReport struct {
	GoVersion string
	GOOS      string
	GOARCH    string
	NumCPU    int
	GOPATH    string
	GOROOT    string
	Compiler  string
	Version   string
}

func Init(out io.Writer, in io.Reader) {
	if in == nil {
		logger.Fatal("the input can't be nil")
	}

	banner, err := ioutil.ReadAll(in)
	if err != nil {
		logger.Fatalf("read the banner error: %s", err)
	}

	show(out, string(banner))
}

func show(out io.Writer, content string) {
	t, err := template.New("banner").Funcs(template.FuncMap{}).Parse(content)
	err = t.Execute(out, RuntimeReport{
		helper.GoVersion(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		hictl.GoPath,
		runtime.GOROOT(),
		runtime.Compiler,
		color.Cyan(hictl.Version),
	})
	if err != nil {
		logger.Error(err.Error())
	}
}
