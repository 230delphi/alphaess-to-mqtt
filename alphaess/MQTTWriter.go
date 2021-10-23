package alphaess

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

// TODO Unit tests for MQTTWriter

// MQTTWriter is a crude implementation to parse the stream for objects, and publish to MQTT
type MQTTWriter struct {
	mySource string
}

func (into *MQTTWriter) Write(newData []byte) (n int, err error) {
	n = len(newData)
	into.parseForJSON(newData)
	return n, err
}

func (into *MQTTWriter) close() (err error) {
	return err
}

func ProcessStream(myReader *bufio.Reader) (found int) {
	var counter = 0 // Read data from stdin in a loop
	//var found int = 0
	var data []byte
	var myWriter = MQTTWriter{}
	var err error
	for {
		_, err = myReader.Peek(MinMessageSize)
		if err == io.EOF && (len(data) < MinMessageSize || bytes.Index(data, []byte("}")) < 0 || bytes.Index(data, []byte("{")) < 0) {
			if len(data) > MinMessageSize {
				DebugLog(string(data))
			}
			break
		} else {
			counter++
			nextLine := make([]byte, 12800)
			var n int
			n, err = myReader.Read(nextLine)
			if err == nil && n > 0 {
				data = bytes.Trim(append(data, nextLine...), "\x00")
			} else if n == 0 && len(data) < 80 {
				DebugLog("EOF error. reads:" + strconv.Itoa(counter) + "data size:" + strconv.Itoa(len(data)))
				if len(data) > MinMessageSize {
					DebugLog(string(data))
				}
				break
			} else {
				ExceptionLog(err, "ProcessStream()")
			}
		}
		myWriter.parseForJSON(data)
		found++
	}
	return found
}

func (into *MQTTWriter) parseForJSON(data []byte) {
	//TODO print warning when we get disconnected:
	// {"MsgType":"Socket","MsgContent":"Socket is closed!","Description":"OK"}
	// and after successful reconnect: client type:ConfigRS : {"SN":"AL2002321010043","Address":.....
	myRecord, _, _, checksumU16 := parseMessage(data)
	var err error
	var start = bytes.Index(myRecord, []byte("{\""))
	var end = bytes.Index(myRecord, []byte("\"}")) + 2
	if ValidateChecksum(data, checksumU16) && start >= 0 && end > 0 {
		var obj Response
		obj, err = UnmarshalJSON(myRecord)
		if obj != nil {
			if publishAlphaESSStats(obj, into.mySource) && !gSystemStarted {
				gSystemStarted = true // only set true after we successfully publish something.
			}
		} else if err != nil {
			ErrorLog("SRC:" + into.mySource + " Unmarshal failed:" + err.Error() + " data:" + string(myRecord))
		}
	}
}
