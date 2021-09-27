package alphaess

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/230delphi/go-any-proxy/anyproxy"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namsral/flag"
	"io"
	"log"
	"strconv"
	"strings"
)

const MinMessageSize = 10

var data = make([]byte, 0)

var gClient mqtt.Client
var gMQTTBrokerAddress string
var gMQTTUser string
var gMQTTPassword string
var gAlphaEssInstance string
var gLogList string
var gLocation string
var gMQTTTimeoutSeconds int
var gTopicBase string
var gMQTTTopic string
var gProxyConnectionImpl string

//var gDebug int

func initFlagConfig() {
	flag.StringVar(&gMQTTBrokerAddress, "MQTTAddress", "tcp://127.0.0.1:1883", "MQTT address. Example: tcp://127.0.0.1:1883\n")
	flag.StringVar(&gMQTTUser, "MQTTUser", "", "MQTT username\n")
	flag.StringVar(&gMQTTPassword, "MQTTPassword", "", "MQTT password\n")
	flag.IntVar(&gMQTTTimeoutSeconds, "MQTTSendTimeout", 5, "MQTT timeout for sending message\n")
	flag.StringVar(&gTopicBase, "MQTTTopicBase", "homeassistant/sensor/", "MQTT base topic. ")
	flag.StringVar(&gAlphaEssInstance, "AlphaESSID", "alphaess1", "AlphaESS instance name, appended to MQTTTopicBase. All data is set on this topic eg: homeassisant/sensor/alphaess1/config\n")
	flag.StringVar(&gLogList, "MSGLogging", "", "Messages to Log. Leave unset for no logging. Log all:\"*\"; log selected: \"GenericRQ,CommandIndexRQ,CommandRQ,ConfigRS,StatusRQ\"")
	flag.StringVar(&gLocation, "TZLocation", "Local", "Timezone override to ensure time of collection is accurate.")
	gMQTTTopic = gTopicBase + gAlphaEssInstance
	DebugLog("initFlagConfig complete")
}

func GetMQTTConnection() (result anyproxy.ProxyConnectionManager) {
	//initFlagConfig() // unnecessary? problem?
	switch gProxyConnectionImpl {
	case "DirectProxyConnection":
		DebugLog(gProxyConnectionImpl)
		result = &anyproxy.DirectProxyConnection{}
	case "LoggingProxyConnection":
		DebugLog(gProxyConnectionImpl)
		result = &anyproxy.LoggingProxyConnection{}
	case "MQTTReadProxyConnection":
		DebugLog(gProxyConnectionImpl)
		result = &MQTTReadProxyConnection{}
	case "MQTTInjectProxyConnection":
		DebugLog(gProxyConnectionImpl)
		result = &MQTTInjectProxyConnection{}
	case "MQTTStubProxyConnection":
		DebugLog(gProxyConnectionImpl)
		result = &MQTTStubProxyConnection{}
	default:
		ErrorLog("ProxyConnectionManager implementation not available. please check configuration: " + gProxyConnectionImpl)
		panic("ProxyConnectionManager implementation not available. please check configuration: " + gProxyConnectionImpl)
	}
	return result
}

type MQTTWriter struct {
	mySource string
}

func (into *MQTTWriter) Write(newData []byte) (n int, err error) {
	DebugLog("MQTTWriter write()" + string(newData))
	data = append(data, newData...)
	into.parseForJSON()
	n = len(newData)
	return n, err
}

func (into *MQTTWriter) close() (err error) {
	return err
}

func ProcessStream(myReader *bufio.Reader) (found int) {
	var counter = 0 // Read data from stdin in a loop
	//var found int = 0
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
				ExceptionLog(err)
			}
		}
		myWriter.parseForJSON()
		found++
	}
	return found
}

func (into *MQTTWriter) parseForJSON() {
	var err error
	var start = bytes.Index(data, []byte("{\""))
	var end = bytes.Index(data, []byte("\"}")) + 2
	var counter = 0
	for start >= 0 && end > MinMessageSize {
		counter++
		data = data[start:]
		end = bytes.Index(data, []byte("\"}")) + 2
		if end > MinMessageSize {
			var myRecord = data[0:end]
			if bytes.Index(myRecord[2:], []byte("{\"")) > 0 {
				ErrorLog("SUSPECT:: more than one" + string(data))
			}
			//trim and discard rest
			data = bytes.Trim(data[(end):], "\x00")
			// deal with record
			var obj Response
			obj, err = UnmarshalJSON(myRecord)
			if obj != nil {
				//DebugLog(fmt.Sprint("success:", obj))
				publishAlphaESSStats(obj, into.mySource)
			} else if err != nil {
				ErrorLog("Unmarshal failed:" + err.Error())
			}
		}
		start = bytes.Index(data, []byte("{\""))
		end = bytes.Index(data, []byte("\"}")) + 2
	}
	if len(data) < MinMessageSize {
		DebugLog("ParseJason() complete: counter:" + strconv.Itoa(counter) + "data:" + string(data))
	}
	if len(data) > 0 && start < 0 {
		// empty buffer of non Json data
		data = make([]byte, 0)
	}
}

func init() {
	initFlagConfig()
	anyproxy.InitConfig()
	gClient = initMQTT()
}

func logResponse(obj Response, source string) {
	//remove package preface "alphaess."
	var myType = strings.ReplaceAll(fmt.Sprintf("%T", obj), "alphaess.", "")
	objString, _ := json.Marshal(obj)
	if strings.Contains(gLogList, myType) {
		log.Println("SRC:" + source + " type:" + myType + " : " + string(objString))
	} else if strings.Index(gLogList, "*") == 0 {
		log.Println("SRC:" + source + " type:" + myType + " : " + string(objString))
	} else {
		DebugLog("SRC:" + source + " type: " + myType + " not in: " + gLogList)
	}
}

func publishAlphaESSStats(obj Response, source string) {
	var destination string
	logResponse(obj, source)
	switch v := obj.(type) {
	case StatusRQ:
		destination = "/state"
		res, _ := json.Marshal(v)
		publishMQTT(gClient, gMQTTTopic+destination, string(res))
		InfoLog(fmt.Sprint("MQTT Published::", v))
	case BatteryRQ:
		destination = "/other"
		res, _ := json.Marshal(v)
		publishMQTT(gClient, gMQTTTopic+destination, string(res))
		InfoLog(fmt.Sprint("MQTT Published::", v))
	case ConfigRS:
		destination = "/attributes"
		res, _ := json.Marshal(v)
		publishMQTT(gClient, gMQTTTopic+destination, string(res))
		InfoLog(fmt.Sprint("MQTT Published::", v))
	default:
		DebugLog(fmt.Sprint("Ignored message type::", v))
	}
}
