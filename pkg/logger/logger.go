package logger

import (
	`errors`
	`fmt`
	`io`
	`os`
	`path/filepath`
	`sync`
	`sync/atomic`
	`text/template`
	`time`

	`github.com/photowey/hictl/pkg/color`
)

const (
	DebugLevel = iota
	ErrorLevel
	FatalLevel
	CriticalLevel
	SuccessLevel
	WarnLevel
	InfoLevel
	HintLevel
)

var (
	sequenceNo uint64
	instance   *Logger
)

var (
	printTemplate  *template.Template
	simpleTemplate *template.Template
	debugTemplate  *template.Template
)

var (
	debugMode = os.Getenv("DEBUG_ENABLED") == "1"
	logLevel  = InfoLevel
)
var (
	errInvalidLogLevel = errors.New("logger: invalid log level")
)

var (
	lock = &sync.Mutex{}
)

type Record struct {
	ID       string
	Level    string
	Message  string
	Filename string
	LineNo   int
}

type Logger struct {
	lock   sync.Mutex
	output io.Writer
}

func (log *Logger) Output(w io.Writer) {
	log.lock.Lock()
	defer log.lock.Unlock()
	log.output = color.NewColorWriter(w)
}
func GetInstance(w io.Writer) *Logger {
	if nil == instance {
		lock.Lock()
		defer lock.Unlock()
		if nil == instance {
			var (
				err          error
				printFormat  = `{{.Message}}{{EndLine}}`
				simpleFormat = `{{Now "2006-01-02 15:04:05.999"}} {{.Level}} ▶ {{.ID}} {{.Message}}{{EndLine}}`
				debugFormat  = `{{Now "2006-01-02 15:04:05.999"}} {{.Level}} ▶ {{.ID}} {{.Filename}}:{{.LineNo}} {{.Message}}{{EndLine}}`
			)

			funcs := template.FuncMap{
				"Now":     Now,
				"EndLine": EndLine,
			}
			printTemplate, err = template.New("print").Funcs(funcs).Parse(printFormat)
			if err != nil {
				panic(err)
			}
			simpleTemplate, err = template.New("simple").Funcs(funcs).Parse(simpleFormat)
			if err != nil {
				panic(err)
			}
			debugTemplate, err = template.New("debug").Funcs(funcs).Parse(debugFormat)
			if err != nil {
				panic(err)
			}

			instance = &Logger{output: color.NewColorWriter(w)}
		}
	}

	return instance
}

func (log *Logger) levelTag(level int) string {
	switch level {
	case FatalLevel:
		return "FATAL"
	case SuccessLevel:
		return "SUCCESS"
	case HintLevel:
		return "HINT"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case CriticalLevel:
		return "CRITICAL"
	default:
		panic(errInvalidLogLevel)
	}
}

func (log *Logger) colorLevel(level int) string {
	switch level {
	case CriticalLevel:
		return color.RedBold(log.levelTag(level))
	case FatalLevel:
		return color.RedBold(log.levelTag(level))
	case InfoLevel:
		return color.BlueBold(log.levelTag(level))
	case HintLevel:
		return color.CyanBold(log.levelTag(level))
	case DebugLevel:
		return color.YellowBold(log.levelTag(level))
	case ErrorLevel:
		return color.RedBold(log.levelTag(level))
	case WarnLevel:
		return color.YellowBold(log.levelTag(level))
	case SuccessLevel:
		return color.GreenBold(log.levelTag(level))
	default:
		panic(errInvalidLogLevel)
	}
}

func (log *Logger) print(message string, args ...any) {
	log.lock.Lock()
	defer log.lock.Unlock()

	record := Record{
		Message: fmt.Sprintf(message, args...),
	}

	err := printTemplate.Execute(log.output, record)
	if err != nil {
		panic(err)
	}
}

func (log *Logger) log(level int, message string, args ...any) {
	if level > logLevel {
		return
	}
	log.lock.Lock()
	defer log.lock.Unlock()

	record := Record{
		ID:      fmt.Sprintf("%04d", atomic.AddUint64(&sequenceNo, 1)),
		Level:   log.colorLevel(level),
		Message: fmt.Sprintf(message, args...),
	}

	err := simpleTemplate.Execute(log.output, record)
	if err != nil {
		panic(err)
	}
}

func (log *Logger) debug(message string, file string, line int, args ...any) {
	if !debugMode {
		return
	}
	log.Output(os.Stderr)

	record := Record{
		ID:       fmt.Sprintf("%04d", atomic.AddUint64(&sequenceNo, 1)),
		Level:    log.colorLevel(DebugLevel),
		Message:  fmt.Sprintf(message, args...),
		LineNo:   line,
		Filename: filepath.Base(file),
	}
	err := debugTemplate.Execute(log.output, record)
	if err != nil {
		panic(err)
	}
}

var log = GetInstance(os.Stdout)

func Printf(message string, args ...any) {
	log.print(message, args...)
}

func Debug(message string, file string, line int) {
	log.debug(message, file, line)
}

func Debugf(message string, file string, line int, args ...any) {
	log.debug(message, file, line, args...)
}

func Info(message string) {
	log.log(InfoLevel, message)
}

func Infof(message string, args ...any) {
	log.log(InfoLevel, message, args...)
}

func Warn(message string) {
	log.log(WarnLevel, message)
}

func Warnf(message string, args ...any) {
	log.log(WarnLevel, message, args...)
}

func Error(message string) {
	log.log(ErrorLevel, message)
}

func Errorf(message string, args ...any) {
	log.log(ErrorLevel, message, args...)
}

func Fatal(message string) {
	log.log(FatalLevel, message)
	os.Exit(255)
}

func Fatalf(message string, args ...any) {
	log.log(FatalLevel, message, args...)
	os.Exit(255)
}

func Success(message string) {
	log.log(SuccessLevel, message)
}

func Successf(message string, args ...any) {
	log.log(SuccessLevel, message, args...)
}

func Hint(message string) {
	log.log(HintLevel, message)
}

func Hintf(message string, args ...any) {
	log.log(HintLevel, message, args...)
}

func Critical(message string) {
	log.log(CriticalLevel, message)
}

func Criticalf(message string, args ...any) {
	log.log(CriticalLevel, message, args...)
}

func Now(layout string) string {
	return time.Now().Format(layout)
}

func EndLine() string {
	return "\n"
}
