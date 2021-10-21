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
	//DebugLog("MQTTWriter write()" + string(newData))
	n = len(newData)
	DebugLog("MQTTWriter write():" + strconv.Itoa(n))
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
	var err error
	var start = bytes.Index(data, []byte("{\""))
	var end = bytes.Index(data, []byte("\"}")) + 2
	var counter = 0
	var header []byte
	var checksum []byte
	for start >= 0 && end > MinMessageSize {
		counter++
		if start > 0 {
			header = data[0:start]
			data = data[start:]
		} else {
			header = nil
			// no change to data
		}
		end = bytes.Index(data, []byte("\"}")) + 2
		if end > MinMessageSize {
			if end < len(data) {
				checksum = data[end:]
			} else {
				checksum = nil
			}
			var myRecord = data[0:end]
			if bytes.Index(myRecord[2:], []byte("{\"")) > 0 {
				ErrorLog("SRC:" + into.mySource + "; SUSPECT:: more than one" + string(data))
			}
			testCheckSum(header, myRecord, checksum)
			var obj Response
			obj, err = UnmarshalJSON(myRecord)
			if obj != nil {
				//DebugLog("SRC:" + into.mySource + "; Obj found; header:'" + string(header) + "'::"
				//+ strconv.Itoa(len(header)) + " checksum:'" + string(checksum) + "'"+ "::"+ strconv.Itoa(len(checksum)))
				publishAlphaESSStats(obj, into.mySource)
			} else if err != nil {
				ErrorLog("SRC:" + into.mySource + " Unmarshal failed:" + err.Error() + " data:" + string(myRecord))
			}
		}
		start = bytes.Index(data, []byte("{\""))
		end = bytes.Index(data, []byte("\"}")) + 2
		if start > 0 {
			data = data[start:]
			DebugLog("this should never happen?")
		} else {
			data = nil
		}
	}
	if len(data) > 0 && start < 0 {
		// empty buffer of non Json data
		data = make([]byte, 0)
	}
}
