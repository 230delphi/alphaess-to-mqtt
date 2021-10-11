package alphaess

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

var gLoggingLevel int = 0 // 0 Error only; 1 Debug; 2 Info

func CheckError(e error) {
	if e != nil {
		DebugLog(e.Error())
		panic(e)
	}
}

func DebugLog(msg string) {
	if gLoggingLevel == 1 {
		log.Println(msg)
	}
}

func InfoLog(msg string) {
	if gLoggingLevel > 0 {
		log.Println(msg)
	}
}

func ErrorLog(msg string) {
	log.Println("ERROR:" + msg)
}

func ExceptionLog(errorResult error) {
	if errorResult != nil {
		ErrorLog("EXP:" + errorResult.Error())
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
