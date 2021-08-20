package alphaess

import (
	"bufio"
	"bytes"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io"
	"strconv"
	"time"
	anyproxy "github.com/230delphi/go-any-proxy/anyproxy"
)

//ProxyConnection "github.com/230delphi/go-any-proxy"
//"../rk-go-any-proxy/main"
// TODO export to config file
// mqtt config
var mqttBrokerAddress = "tcp://MASK:1883"
var mqttUser = "mqttMASK"
var mqttPassword = "MASK"
var gclient mqtt.Client
var alphaEssInstance = "alphaess1"

// const TopicBase = "alpha-ess/sensor/alpha-ess/"
const TopicBase = "homeassistant/sensor/"

var MQTTTopic = TopicBase + alphaEssInstance
var myProxyConnection ProxyConnection

type HasMQTTConfig struct {
	DeviceClass       string `json:"device_class"`
	Name              string `json:"name"`
	StateTopic        string `json:"state_topic"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	ValueTemplate     string `json:"value_template"`
}

func getUniqueFilename(prepend string) string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10) + "_" + prepend + "_stream.file"
}

func processStream(err error, myreader *bufio.Reader) {
	data := make([]byte, 0)
	var counter = 0 // Read data from stdin in a loop

	for {
		debugLog("left:" + strconv.Itoa(len(data)))
		_, err = myreader.Peek(10)
		//debugLog("peek done" )
		if err == io.EOF && (len(data) < 80 || bytes.Index(data, []byte("\"}")) < 0) {
			debugLog("peek done - EOF. reads:" + strconv.Itoa(counter))
			break
		} else {
			//if len(data)==0 || start == -1
			//debugLog("reading. current buffer:" + strconv.Itoa(len(data)) + ":" + strconv.Itoa(start) )
			counter++
			nextLine := make([]byte, 12800)
			_, err = myreader.Read(nextLine)
			data = append(data, nextLine...)
		}
		start := bytes.Index(data, []byte("{\""))
		//debugLog("found:" + strconv.Itoa(start) +":" + string(data))

		if start >= 0 {
			data = data[start:]
			//start = bytes.Index(data, []byte("{\""))
			end := bytes.Index(data, []byte("\"}")) + 2
			//debugLog("found end:" + strconv.Itoa(end) +":" + string(data))
			if end > 10 {
				var myrecord = data[0:end]
				//var newstart = bytes.Index(myrecord[1:], []byte("{\""))
				//if newstart < end{
				//	var discarded []byte = myrecord [:newstart]
				//	debugLog("DISCARDED------ unterminated: " +string(discarded))
				//	myrecord = myrecord[newstart:]
				if bytes.Index(myrecord[2:], []byte("{\"")) > 0 {
					debugLog("SUSPECT::" + string(data))
				}
				//trim and discard rest
				data = bytes.Trim(data[(end):], "\x00")
				//myrecord = bytes.Trim(myrecord, "\x00")
				debugLog("RECORD::" + string(myrecord[len(myrecord)-1]) + "::" + string(myrecord))
				// deal with record
				var obj AlphaessResponse
				obj, err = UnmarshalJSON(([]byte)(myrecord))
				if obj != nil {
					//an, _ := json.Marshal(obj)
					debugLog(fmt.Sprint("success:", obj))
					publishAlphaESSStats(obj, gclient)
				} else if err != nil {
					errorLog("failed:" + err.Error())
				}
			}
		}
	}
}

func init() {
	gclient = initMQTT()
}
