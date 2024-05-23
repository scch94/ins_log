package ins_log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

var logger = log.New(os.Stdout, "", 0)
var service = "not_specified"
var logLevel = 6

// Clave para almacenar y recuperar el UTFI en el contexto
type contextKey string

const packageKey = contextKey("packageName")
const utfIKey = contextKey("UTFI")

func StartLogger() {
	logger.SetOutput(new(logWriter))
}

func StartLoggerWithWriter(w io.Writer) {
	logger.SetOutput(w)
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

func GenerateUTFI() string {
	return uuid.New().String()[24:]
}

func SetService(name string) {
	service = name
}

func SetLevel(level string) {
	logLevel = levelToInt(level)
}

func Fatal(c context.Context, msg string) {
	doLog(c, 1, msg)
}
func Error(c context.Context, msg string) {
	doLog(c, 2, msg)
}
func Warn(c context.Context, msg string) {
	doLog(c, 3, msg)
}
func Info(c context.Context, msg string) {
	doLog(c, 4, msg)
}
func Debug(c context.Context, msg string) {
	doLog(c, 5, msg)
}
func Trace(c context.Context, msg string) {
	doLog(c, 6, msg)
}
func Print(c context.Context, msg string) {
	doLog(c, 5, msg)
}
func Fatalf(c context.Context, msg string, args ...interface{}) {
	doLog(c, 1, msg, args...)
}
func Errorf(c context.Context, msg string, args ...interface{}) {
	doLog(c, 2, msg, args...)
}
func Warnf(c context.Context, msg string, args ...interface{}) {
	doLog(c, 3, msg, args...)
}
func Infof(c context.Context, msg string, args ...interface{}) {
	doLog(c, 4, msg, args...)
}
func Debugf(c context.Context, msg string, args ...interface{}) {
	doLog(c, 5, msg, args...)
}
func Tracef(c context.Context, msg string, args ...interface{}) {
	doLog(c, 6, msg, args...)
}
func Printf(c context.Context, msg string, args ...interface{}) {
	doLog(c, 5, msg, args...)
}

func doLog(c context.Context, lineLevel int, msg string, params ...interface{}) {
	if lineLevel > logLevel {
		return
	}

	dateTime := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := levelToString(lineLevel)
	utfi := GetUTFIFromContext(c)
	msg = replaceCharacters(msg)
	packageName := GetPackageNameFromContext(c)
	line := fmt.Sprintf("[%s] [%s] [%s] [%s] [%s] %s", dateTime, service, packageName, levelStr, utfi, msg)
	if len(params) == 0 {
		logger.Print(line)
	} else {
		params = mapParams(params)
		logger.Printf(line, params...)
	}
}

func mapParams(params []interface{}) []interface{} {
	mappedParams := []interface{}{}
	for _, i := range params {
		switch v := i.(type) {
		case string:
			mappedParams = append(mappedParams, replaceCharacters(v))
		default:
			mappedParams = append(mappedParams, v)
		}
	}
	return mappedParams
}

func replaceCharacters(s string) string {
	return strings.ReplaceAll(s, "\n", `\n`)
}

func emptyStringIfNil(data interface{}) interface{} {
	if data == nil {
		return ""
	} else {
		return data
	}
}

func levelToString(level int) string {
	switch level {
	case 0:
		return "none"
	case 1:
		return "fatal"
	case 2:
		return "error"
	case 3:
		return "warn"
	case 4:
		return "info"
	case 5:
		return "debug"
	default:
		return "trace"
	}
}

func levelToInt(level string) int {
	switch level {
	case "none":
		return 0
	case "fatal":
		return 1
	case "error":
		return 2
	case "warn":
		return 3
	case "info":
		return 4
	case "debug":
		return 5
	default:
		return 6
	}
}

// Helper para agregar y obtener UTFI del contexto
func SetUTFIInContext(ctx context.Context, utfI string) context.Context {
	return context.WithValue(ctx, utfIKey, utfI)
}

func GetUTFIFromContext(c context.Context) string {
	if utfI, ok := c.Value(utfIKey).(string); ok {
		return utfI
	}
	return ""
}

// Helper para agregar y obtener packageName del contexto
func SetPackageNameInContext(ctx context.Context, packageName string) context.Context {
	return context.WithValue(ctx, packageKey, packageName)
}

func GetPackageNameFromContext(c context.Context) string {
	if packageName, ok := c.Value(packageKey).(string); ok {
		return packageName
	}
	return ""
}
