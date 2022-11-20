package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"

	// go.etcd.io/etcd imports capnslog, which calls log.SetOutput in its
	// init() function, so importing it here means that our log.SetOutput
	// wins. this is fixed in coreos v3.5, which is not released yet. See
	// https://github.com/etcd-io/etcd/issues/12498 for more information.
	_ "github.com/coreos/pkg/capnslog"
	"github.com/hashicorp/go-hclog"
	v1 "github.com/neermitt/opsos/api/v1"
)

var (
	// ValidLevels are the log level names that Terraform recognizes.
	ValidLevels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "OFF"}

	// logger is the global hclog logger
	logger hclog.Logger

	// logWriter is a global writer for logs, to be used with the std log package
	logWriter io.Writer
)

func init() {
	logger = hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:              "",
		Level:             hclog.DefaultLevel,
		Output:            io.Writer(os.Stdout),
		IndependentLevels: true,
		JSONFormat:        false,
	})
}

func InitLogger(conf v1.ConfigSpec) {
	logger = newHCLogger("", conf.Logs)
	logWriter = logger.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})

	// set up the default std library logger to use our output
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(logWriter)
}

// RegisterSink adds a new log sink which writes all logs to the given file.
func RegisterSink(f *os.File) {
	l, ok := logger.(hclog.InterceptLogger)
	if !ok {
		panic("global logger is not an InterceptLogger")
	}

	if f == nil {
		return
	}

	l.RegisterSink(hclog.NewSinkAdapter(&hclog.LoggerOptions{
		Level:  hclog.Trace,
		Output: f,
	}))
}

// LogOutput return the default global log io.Writer
func LogOutput() io.Writer {
	return logWriter
}

// HCLogger returns the default global hclog logger
func HCLogger() hclog.Logger {
	return logger
}

// newHCLogger returns a new hclog.Logger instance with the given name
func newHCLogger(name string, logConf v1.LogSpec) hclog.Logger {
	logOutput := io.Writer(os.Stderr)
	logLevel, json := globalLogLevel(logConf)

	if logPath := logConf.File; logPath != nil {
		f, err := os.OpenFile(*logPath, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		} else {
			logOutput = f
		}
	}

	return hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:              name,
		Level:             logLevel,
		Output:            logOutput,
		IndependentLevels: true,
		JSONFormat:        json,
	})
}

// CurrentLogLevel returns the current log level string based the environment vars
func CurrentLogLevel(logConf v1.LogSpec) string {
	ll, _ := globalLogLevel(logConf)
	return strings.ToUpper(ll.String())
}

func globalLogLevel(logConf v1.LogSpec) (hclog.Level, bool) {
	logLevel := ""
	if logConf.Level != nil {
		logLevel = *logConf.Level
	}
	return parseLogLevel(strings.ToUpper(logLevel)), logConf.JSON
}

func parseLogLevel(envLevel string) hclog.Level {
	if envLevel == "" {
		return hclog.Off
	}
	if envLevel == "JSON" {
		envLevel = "TRACE"
	}

	logLevel := hclog.Trace
	if isValidLogLevel(envLevel) {
		logLevel = hclog.LevelFromString(envLevel)
	} else {
		fmt.Fprintf(os.Stderr, "[WARN] Invalid log level: %q. Defaulting to level: TRACE. Valid levels are: %+v",
			envLevel, ValidLevels)
	}

	return logLevel
}

// IsDebugOrHigher returns whether the current log level is debug or trace
func IsDebugOrHigher(logConf v1.LogSpec) bool {
	level, _ := globalLogLevel(logConf)
	return level == hclog.Debug || level == hclog.Trace
}

func isValidLogLevel(level string) bool {
	for _, l := range ValidLevels {
		if level == string(l) {
			return true
		}
	}

	return false
}
