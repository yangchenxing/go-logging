package logging

import (
	"net"
	"os"
	"runtime"
	"strings"
)

var (
	packageHosts = []string{
		"github.com",
		"bitbucket.org",
		"launchpad.net",
		"google.golang.org",
		"golang.org",
		"gopkg.in",
	}
	globalStaticContext map[string]interface{}
	handlers            = make(map[string][]*Handler)
	localContextCache   = make(map[uintptr]map[string]interface{})
	unknownLocalContext = map[string]interface{}{
		"file": "",
		"line": 0,
		"func": "",
	}
	timeFormat = "2006-01-02:15:04:05-0700"
)

// AddPackageHost add a customized golang package host to known package hosts.
// While create static local context for logging message formatting, known package hosts is used
// to cut off GOPATH directory from the path of file. It ensures output compiled in different GOPATH
// will log with same path of file.
func AddPackageHost(host string) {
	packageHosts = append(packageHosts, host)
}

// SetTimeFormat set the time format.
func SetTimeFormat(format string) {
	timeFormat = format
}

func init() {
	// 获取本地IP
	ip := ""
	if infs, err := net.Interfaces(); err == nil && len(infs) > 0 {
	InfLoop:
		for _, inf := range infs {
			if inf.Flags&net.FlagLoopback != 0 {
				continue
			}
			if addrs, err := inf.Addrs(); err == nil && len(addrs) > 0 {
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok {
						if ipv4 := ipnet.IP.To4(); ipv4 != nil {
							ip = ipv4.String()
							break InfLoop
						}
					}
				}
			}
		}
	}
	hostname, _ := os.Hostname()
	// 构造全局静态上下文信息
	globalStaticContext = map[string]interface{}{
		"ip":       ip,
		"hostname": hostname,
	}
	// 构造默认handler
	defaultHandler := &Handler{
		Levels:  []string{"debug", "info", "warn", "error", "fatal"},
		Format:  "$level [$time][$file:$line][$func] $message",
		Writers: []Writer{os.Stderr},
	}
	for _, level := range defaultHandler.Levels {
		handlers[level] = []*Handler{defaultHandler}
	}
}

func getLocalStaticContext(skip int) map[string]interface{} {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return unknownLocalContext
	}
	context := localContextCache[pc]
	if context != nil {
		return context
	}
	funcname := ""
	caller := runtime.FuncForPC(pc)
	if caller != nil {
		funcname = caller.Name()
	}
	for _, host := range packageHosts {
		pos := strings.Index(file, host)
		if pos >= 0 {
			file = file[pos:]
		}
	}
	context = map[string]interface{}{
		"file": file,
		"line": line,
		"func": funcname,
	}
	localContextCache[pc] = context
	return context
}

// LogWithSkip log message with specified skip
func LogWithSkip(skip int, level string, format string, args ...interface{}) {
	defaultSession.LogWithSkip(skip+1, level, format, args...)
}

// Log log message with default session.
func Log(level, format string, args ...interface{}) {
	defaultSession.LogWithSkip(3, level, format, args...)
}

// Debug log message with "debug" level and default session.
func Debug(format string, args ...interface{}) {
	defaultSession.LogWithSkip(3, "debug", format, args...)
}

// Info log message with "info" level and default session.
func Info(format string, args ...interface{}) {
	defaultSession.LogWithSkip(3, "info", format, args...)
}

// Warn log message with "warn" level and default session.
func Warn(format string, args ...interface{}) {
	defaultSession.LogWithSkip(3, "warn", format, args...)
}

// Error log message with "error" level and default session.
func Error(format string, args ...interface{}) {
	defaultSession.LogWithSkip(3, "error", format, args...)
}

// Fatal log message with "fatal" level and default session.
func Fatal(format string, args ...interface{}) {
	defaultSession.LogWithSkip(3, "fatal", format, args...)
}
