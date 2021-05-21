package pretty

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"github.com/gen0cide/genolog"
)

var (
	// Logger is a global singleton logger
	Logger genolog.Logger
)

var (
	defaultProg  = `APP`
	startName    = `cli`
	defaultLevel = logrus.InfoLevel

	global *logrus.Logger
)

type prettyLogger struct {
	internal *logrus.Logger
	writer   *logWriter
	prog     string
	context  string
}

type logWriter struct {
	Name   string
	Prog   string
	Output io.Writer
}

// NewPrettyLogger returns a new pretty console logger implementing the Logger interface.
func NewPrettyLogger(prog, name string, existing genolog.Logger) genolog.Logger {
	if prog == "" {
		prog = defaultProg
	}

	if name == "" {
		name = startName
	}

	if existing != nil {
		existing.SetName(name)
		existing.SetProg(prog)
		return existing
	}

	if global == nil {
		global = logrus.New()
		global.SetLevel(defaultLevel)
	}

	writer := &logWriter{
		Name:   name,
		Prog:   prog,
		Output: color.Output,
	}
	global.Out = writer

	logger := &prettyLogger{
		internal: global,
		writer:   writer,
		prog:     prog,
		context:  name,
	}

	global.Formatter = logger

	Logger = logger

	return logger
}

// ZapLogger implements the genolog.Logger interface.
func (p *prettyLogger) ZapLogger() *zap.SugaredLogger {
	return nil
}

// GetWriter implements the genolog.Logger interface.
func (p *prettyLogger) GetWriter() io.Writer {
	return p.writer
}

// LogrusLogger implements the genolog.Logger interface.
func (p *prettyLogger) LogrusLogger() *logrus.Logger {
	return p.internal
}

// SetName implements the genolog.Logger interface.
func (p *prettyLogger) SetName(s string) {
	p.context = s
	p.writer.Name = s
}

// SetProg implements the genolog.Logger interface.
func (p *prettyLogger) SetProg(s string) {
	p.prog = s
	p.writer.Prog = s
}

// GetOutput implements the genolog.Logger interface.
func (p *prettyLogger) GetOutput() io.Writer {
	return p.writer.Output
}

// SetOutput implements the genolog.Logger interface.
func (p *prettyLogger) SetOutput(w io.Writer) {
	p.writer.Output = w
}

// Format is called to pretty format various output types.
func (p *prettyLogger) Format(entry *logrus.Entry) ([]byte, error) {
	var logLvl, logLine string

	buf := new(bytes.Buffer)

	if len(entry.Data) > 0 {
		buf.WriteString(genolog.NoColor(" "))
		buf.WriteString("\n")
		buf.WriteString(genolog.BoldBrightWhite(">>>"))
		buf.WriteString(" ")
		switch entry.Level.String() {
		case "debug":
			buf.WriteString(genolog.BoldCyan("Debug Details"))
		case "info":
			buf.WriteString(genolog.BoldBrightWhite("Info Details"))
		case "warning":
			buf.WriteString(genolog.BoldBrightBlue("Warning Details"))
		case "error":
			buf.WriteString(genolog.BoldBrightYellow("Error Details"))
		case "fatal":
			buf.WriteString(genolog.BoldBrightRed("Fatal Details"))
		default:
			buf.WriteString(genolog.BoldBrightGreen("Details"))
		}
		buf.WriteString(genolog.NoColor("\n"))
		for k, v := range entry.Data {
			buf.WriteString(
				fmt.Sprintf(
					"%20s %s %s\n",
					genolog.BoldYellow("%20s", k),
					genolog.NoColor("="),
					genolog.BrightWhite("%v", v),
				),
			)
		}
	}
	switch entry.Level.String() {
	case "debug":
		logLvl = genolog.BoldCyan("DEBUG")
		logLine = genolog.Cyan(entry.Message)
	case "info":
		logLvl = genolog.BoldBrightWhite(" INFO")
		logLine = genolog.BrightWhite(entry.Message)
	case "warning":
		logLvl = genolog.BoldBrightBlue(" WARN")
		logLine = genolog.BrightBlue(entry.Message)
	case "error":
		logLvl = genolog.BoldBrightYellow("ERROR")
		logLine = genolog.BrightYellow(entry.Message)
	case "fatal":
		logLvl = genolog.BoldBrightRed("FATAL")
		logLine = genolog.BrightYellow(entry.Message)
	default:
		logLvl = genolog.BoldBrightGreen(" LOG ")
		logLine = genolog.Green(entry.Message)
	}
	line := fmt.Sprintf(
		"%s %s%s\n",
		logLvl,
		logLine,
		buf.String(),
	)
	return []byte(line), nil
}

// Print implements the genolog.Logger interface.
func (p *prettyLogger) Print(args ...interface{}) {
	p.internal.Print(args...)
	return
}

// Printf implements the genolog.Logger interface.
func (p *prettyLogger) Printf(format string, args ...interface{}) {
	p.internal.Printf(format, args...)
	return
}

// Println implements the genolog.Logger interface.
func (p *prettyLogger) Println(args ...interface{}) {
	p.internal.Println(args...)
	return
}

// Printw implements the genolog.Logger interface.
func (p *prettyLogger) Printw(body string, args ...interface{}) {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i++ {
		if i%2 == 0 || i == 0 {
			fields[fmt.Sprintf("%v", args[i])] = nil
			continue
		}
		fields[fmt.Sprintf("%v", args[i-1])] = args[i]
	}

	p.internal.WithFields(fields).Println(body)
	return
}

// Debug implements the genolog.Logger interface.
func (p *prettyLogger) Debug(args ...interface{}) {
	p.internal.Debug(args...)
	return
}

// Debugf implements the genolog.Logger interface.
func (p *prettyLogger) Debugf(format string, args ...interface{}) {
	p.internal.Debugf(format, args...)
	return
}

// Debugln implements the genolog.Logger interface.
func (p *prettyLogger) Debugln(args ...interface{}) {
	p.internal.Debugln(args...)
	return
}

// Debugw implements the genolog.Logger interface.
func (p *prettyLogger) Debugw(body string, args ...interface{}) {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i++ {
		if i%2 == 0 || i == 0 {
			fields[fmt.Sprintf("%v", args[i])] = nil
			continue
		}
		fields[fmt.Sprintf("%v", args[i-1])] = args[i]
	}

	p.internal.WithFields(fields).Debug(body)
	return
}

// Info implements the genolog.Logger interface.
func (p *prettyLogger) Info(args ...interface{}) {
	p.internal.Info(args...)
	return
}

// Infof implements the genolog.Logger interface.
func (p *prettyLogger) Infof(format string, args ...interface{}) {
	p.internal.Infof(format, args...)
	return
}

// Infoln implements the genolog.Logger interface.
func (p *prettyLogger) Infoln(args ...interface{}) {
	p.internal.Infoln(args...)
	return
}

// Infow implements the genolog.Logger interface.
func (p *prettyLogger) Infow(body string, args ...interface{}) {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i++ {
		if i%2 == 0 || i == 0 {
			fields[fmt.Sprintf("%v", args[i])] = nil
			continue
		}
		fields[fmt.Sprintf("%v", args[i-1])] = args[i]
	}

	p.internal.WithFields(fields).Info(body)
	return
}

// Warn implements the genolog.Logger interface.
func (p *prettyLogger) Warn(args ...interface{}) {
	p.internal.Warn(args...)
	return
}

// Warnf implements the genolog.Logger interface.
func (p *prettyLogger) Warnf(format string, args ...interface{}) {
	p.internal.Warnf(format, args...)
	return
}

// Warnln implements the genolog.Logger interface.
func (p *prettyLogger) Warnln(args ...interface{}) {
	p.internal.Warnln(args...)
	return
}

// Warnw implements the genolog.Logger interface.
func (p *prettyLogger) Warnw(body string, args ...interface{}) {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i++ {
		if i%2 == 0 || i == 0 {
			fields[fmt.Sprintf("%v", args[i])] = nil
			continue
		}
		fields[fmt.Sprintf("%v", args[i-1])] = args[i]
	}

	p.internal.WithFields(fields).Warn(body)
	return
}

// Error implements the genolog.Logger interface.
func (p *prettyLogger) Error(args ...interface{}) {
	p.internal.Error(args...)
	return
}

// Errorf implements the genolog.Logger interface.
func (p *prettyLogger) Errorf(format string, args ...interface{}) {
	p.internal.Errorf(format, args...)
	return
}

// Errorln implements the genolog.Logger interface.
func (p *prettyLogger) Errorln(args ...interface{}) {
	p.internal.Errorln(args...)
	return
}

// Errorw implements the genolog.Logger interface.
func (p *prettyLogger) Errorw(body string, args ...interface{}) {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i++ {
		if i%2 == 0 || i == 0 {
			fields[fmt.Sprintf("%v", args[i])] = nil
			continue
		}
		fields[fmt.Sprintf("%v", args[i-1])] = args[i]
	}

	p.internal.WithFields(fields).Error(body)
	return
}

// Fatal implements the genolog.Logger interface.
func (p *prettyLogger) Fatal(args ...interface{}) {
	p.internal.Fatal(args...)
	return
}

// Fatalf implements the genolog.Logger interface.
func (p *prettyLogger) Fatalf(format string, args ...interface{}) {
	p.internal.Fatalf(format, args...)
	return
}

// Fatalln implements the genolog.Logger interface.
func (p *prettyLogger) Fatalln(args ...interface{}) {
	p.internal.Fatalln(args...)
	return
}

// Fatalw implements the genolog.Logger interface.
func (p *prettyLogger) Fatalw(body string, args ...interface{}) {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i++ {
		if i%2 == 0 || i == 0 {
			fields[fmt.Sprintf("%v", args[i])] = nil
			continue
		}
		fields[fmt.Sprintf("%v", args[i-1])] = args[i]
	}

	p.internal.WithFields(fields).Fatal(body)
	return
}

// Raw implements the genolog.Logger interface.
func (p *prettyLogger) Raw(args ...interface{}) {
	_, _ = fmt.Fprint(p.writer.Output, args...)
	return
}

// Rawf implements the genolog.Logger interface.
func (p *prettyLogger) Rawf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(p.writer.Output, format, args...)
	return
}

// Rawln implements the genolog.Logger interface.
func (p *prettyLogger) Rawln(args ...interface{}) {
	_, _ = fmt.Fprintln(p.writer.Output, args...)
	return
}

// SetLogLevel implements the genolog.Logger interface.
func (p *prettyLogger) SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		p.internal.SetLevel(logrus.DebugLevel)
	case "info":
		p.internal.SetLevel(logrus.InfoLevel)
	case "warn":
		p.internal.SetLevel(logrus.WarnLevel)
	case "error":
		p.internal.SetLevel(logrus.ErrorLevel)
	case "fatal":
		p.internal.SetLevel(logrus.FatalLevel)
	}
}

func (w *logWriter) Write(p []byte) (int, error) {
	output := fmt.Sprintf(
		"%s%s%s%s%s %s",
		genolog.BrightWhite("["),
		genolog.BoldBrightCyan(w.Prog),
		genolog.BrightWhite(":"),
		genolog.BrightGreen(strings.ToLower(w.Name)),
		genolog.BrightWhite("]"),
		string(p),
	)
	written, err := io.Copy(w.Output, strings.NewReader(output))
	return int(written), err
}
