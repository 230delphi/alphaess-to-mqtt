package alphaess

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"testing"
	"time"
)

var gLoggingLevel = 0 // 0 Error only; 1 Debug; 2 Info

func DebugEnabled() bool {
	if gLoggingLevel == 1 {
		return true
	} else {
		return false
	}
}
func InfoEnabled() bool {
	if gLoggingLevel > 0 {
		return true
	} else {
		return false
	}
}

func CheckError(e error) {
	if e != nil {
		DebugLog(e.Error())
		panic(e)
	}
}

func DebugLog(msg string, a ...interface{}) {
	if DebugEnabled() {
		log.Debug(fmt.Sprintf(msg, a...))
	}
}

func InfoLog(msg string, a ...interface{}) {
	if InfoEnabled() {
		log.Infof(msg, a...)
	}
}

func ErrorLog(msg string, a ...interface{}) {
	log.Errorf(msg, a...)
}

func WarningLog(msg string, a ...interface{}) {
	log.Warningf(msg, a...)
}

func ExceptionLog(errorResult error, context string) {
	if errorResult != nil {
		ErrorLog("EXP:" + context + ": " + errorResult.Error())
	}
}

func getUniqueFilename(prepend string) string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10) + "_" + prepend + "_stream.file"
}

func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func AssertNotNil(t *testing.T, a interface{}, message string) {
	if a != nil {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("AssertNotNil %v", a)
	}
	t.Fatal(message)
}
