package genolog

import (
	"io"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// Logger is an interface that allows for generic logging capabilities to be defined
// for the applications that need pluggable logging in various forms.
type Logger interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Printw(body string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Debugw(body string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Infow(body string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	Warnw(body string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Errorw(body string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Fatalw(body string, args ...interface{})

	Raw(args ...interface{})
	Rawf(format string, args ...interface{})
	Rawln(args ...interface{})

	SetLogLevel(s string)
	LogrusLogger() *logrus.Logger
	ZapLogger() *zap.SugaredLogger
	SetName(s string)
	SetProg(s string)
	GetWriter() io.Writer
	GetOutput() io.Writer
	SetOutput(io.Writer)
}

var (
	// BoldBrightGreen defines a color printer
	BoldBrightGreen = color.New(color.FgHiGreen, color.Bold).SprintfFunc()

	// BoldBrightWhite defines a color printer
	BoldBrightWhite = color.New(color.FgHiWhite, color.Bold).SprintfFunc()

	// BoldBrightRed defines a color printer
	BoldBrightRed = color.New(color.FgHiRed, color.Bold).SprintfFunc()

	// BoldBrightYellow defines a color printer
	BoldBrightYellow = color.New(color.FgHiYellow, color.Bold).SprintfFunc()

	// BoldBrightCyan defines a color printer
	BoldBrightCyan = color.New(color.FgHiCyan, color.Bold).SprintfFunc()

	// BoldBrightBlue defines a color printer
	BoldBrightBlue = color.New(color.FgHiBlue, color.Bold).SprintfFunc()

	// BoldBrightMagenta defines a color printer
	BoldBrightMagenta = color.New(color.FgHiMagenta, color.Bold).SprintfFunc()

	// BrightGreen defines a color printer
	BrightGreen = color.New(color.FgHiGreen).SprintfFunc()

	// BrightWhite defines a color printer
	BrightWhite = color.New(color.FgHiWhite).SprintfFunc()

	// BrightRed defines a color printer
	BrightRed = color.New(color.FgHiRed).SprintfFunc()

	// BrightYellow defines a color printer
	BrightYellow = color.New(color.FgHiYellow).SprintfFunc()

	// BrightCyan defines a color printer
	BrightCyan = color.New(color.FgHiCyan).SprintfFunc()

	// BrightBlue defines a color printer
	BrightBlue = color.New(color.FgHiBlue).SprintfFunc()

	// BrightMagenta defines a color printer
	BrightMagenta = color.New(color.FgHiMagenta).SprintfFunc()

	// BoldGreen defines a color printer
	BoldGreen = color.New(color.FgGreen, color.Bold).SprintfFunc()

	// BoldWhite defines a color printer
	BoldWhite = color.New(color.FgWhite, color.Bold).SprintfFunc()

	// BoldRed defines a color printer
	BoldRed = color.New(color.FgRed, color.Bold).SprintfFunc()

	// BoldYellow defines a color printer
	BoldYellow = color.New(color.FgYellow, color.Bold).SprintfFunc()

	// BoldCyan defines a color printer
	BoldCyan = color.New(color.FgCyan, color.Bold).SprintfFunc()

	// BoldBlue defines a color printer
	BoldBlue = color.New(color.FgBlue, color.Bold).SprintfFunc()

	// BoldMagenta defines a color printer
	BoldMagenta = color.New(color.FgMagenta, color.Bold).SprintfFunc()

	// Green defines a color printer
	Green = color.New(color.FgGreen).SprintfFunc()

	// White defines a color printer
	White = color.New(color.FgWhite).SprintfFunc()

	// Red defines a color printer
	Red = color.New(color.FgRed).SprintfFunc()

	// Yellow defines a color printer
	Yellow = color.New(color.FgYellow).SprintfFunc()

	// Cyan defines a color printer
	Cyan = color.New(color.FgCyan).SprintfFunc()

	// Blue defines a color printer
	Blue = color.New(color.FgBlue).SprintfFunc()

	// Magenta defines a color printer
	Magenta = color.New(color.FgMagenta).SprintfFunc()

	// NoColor defines a color printer
	NoColor = color.New(color.Reset).SprintfFunc()
)
