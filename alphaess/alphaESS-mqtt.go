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
	"os"
	"strconv"
	"strings"
	"time"
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
var gChargeBatteryTopic = "/action/chargebattery"
var gAttributesTopic = "/attributes"
var gLastConfigRS ConfigRS
var ClientInject []byte

type ChargeBatteryAction struct {
	Action          string `json:"action,omitempty"`   //startCharge|stopCharge
	Hour            int    `json:"hour,omitempty"`     // optional hour to start
	MinimumDuration int    `json:"duration,omitempty"` // min minutes to be operating
	BatHighCap      int    `json:"unit_of_measurement,omitempty"`
}

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
	DebugLog("initFlagConfig complete" + gMQTTBrokerAddress)
}

func printConfig() {
	fmt.Printf(gMQTTBrokerAddress)
	fmt.Printf(gMQTTUser)

}

func GetMQTTConnection() (result anyproxy.ProxyConnectionManager) {
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

func initMQTT() (myClient mqtt.Client) {
	DebugLog("init MQTT: " + gMQTTBrokerAddress)
	//mqtt.DEBUG = 	log.New(os.Stderr, "[DEBUG][MQTT]", log.Ltime)
	mqtt.WARN = log.New(os.Stderr, "[WARN] [MQTT]", log.Ltime)
	mqtt.ERROR = log.New(os.Stdout, "[ERROR][MQTT]", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] [MQTT]", 0)
	opts := mqtt.NewClientOptions().AddBroker(gMQTTBrokerAddress)
	var secs = time.Now().Unix()
	opts.SetClientID("alphaESSGoClient_" + strconv.Itoa(int(secs))[6:])
	opts.SetUsername(gMQTTUser)
	opts.SetPassword(gMQTTPassword)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	// per common problems, avoids possible deadlock.
	opts.SetOrderMatters(false)
	opts.SetKeepAlive(time.Duration(30) * time.Second)
	opts.SetCleanSession(true)
	opts.SetProtocolVersion(4)
	opts.SetConnectRetryInterval(time.Duration(30) * time.Second)
	opts.SetConnectTimeout(time.Duration(gMQTTTimeoutSeconds) * time.Second)
	opts.SetMaxReconnectInterval(time.Duration(60) * time.Second)
	opts.SetPingTimeout(time.Duration(1) * time.Second)
	opts.SetWriteTimeout(time.Duration(gMQTTTimeoutSeconds) * time.Second)
	opts.SetDefaultPublishHandler(mqtPublishHandler)
	opts.SetConnectionLostHandler(connLostHandler)
	myClient = mqtt.NewClient(opts)
	if !connectClient(myClient) {
		panic("Failed to connect to MQTT: " + gMQTTBrokerAddress)
	}
	gChargeBatteryTopic = gMQTTTopic + gChargeBatteryTopic
	return myClient
}

var mqtPublishHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	DebugLog("TOPIC: " + msg.Topic())
	DebugLog("MSG: " + string(msg.Payload()))
	//if gChargeBatteryTopic
	//gLastConfigRS
	switch msg.Topic() {
	case gChargeBatteryTopic:
		InfoLog("MQTT Received on topic " + msg.Topic())
		var startHour = 0
		var endHour = 0
		//Default to charge for this hour, minimum 30mins, charging stops % default
		//TODO check payload for action (startCharge|stopCharge), optional configuration of hour(00-23), minimum duration (0-300) and % charge (BatHighCap:0-100)
		//ChargeBatteryAction
		t := time.Now()
		startHour = t.Hour()
		if t.Minute() < 30 {
			endHour = startHour + 1
		} else {
			endHour = startHour + 2
		}
		gLastConfigRS.TimeChaF2 = string(startHour)
		gLastConfigRS.TimeChaE2 = string(endHour)
		injectConfig, _ := json.Marshal(gLastConfigRS)
		ClientInject = injectConfig
	default:
		InfoLog("Ignored message Received on topic::" + msg.Topic())
	}
}

func connLostHandler(client mqtt.Client, err error) {
	fmt.Printf("Connection lost, reason: %v\n", err)
	connectClient(client)
	ErrorLog("Reconnected MQTT")
}

func connectClient(myClient mqtt.Client) (success bool) {
	if token := myClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if myClient.IsConnectionOpen() {
		DebugLog("initMQTT: connected")
		success = true
	} else {
		ErrorLog("Failed to connect to MQTT: " + gMQTTBrokerAddress)
		return false
	}
	return success
}

func publishMQTT(mqClient mqtt.Client, topic string, msg string) {
	t := mqClient.Publish(topic, 0, true, msg)
	go func() {
		var success = t.WaitTimeout(time.Duration(gMQTTTimeoutSeconds) * time.Second)
		if t.Error() != nil {
			fmt.Printf("Failed to send message: %v\n", t.Error())
			ExceptionLog(t.Error())
		}
		if !success {
			ErrorLog("Timeout occurred after: " + strconv.Itoa(gMQTTTimeoutSeconds) + " sending:" + msg)
			ExceptionLog(t.Error())
		}
	}()
}

func subscribeTopic(mqClient mqtt.Client, topic string) {
	token := mqClient.Subscribe(topic, 1, nil)
	token.Wait()
	DebugLog("Subscribed to topic: " + topic)
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
	var header string
	var checksum string
	for start >= 0 && end > MinMessageSize {
		counter++
		header = string(data[0:(start - 1)])
		data = data[start:]
		end = bytes.Index(data, []byte("\"}")) + 2
		if end > MinMessageSize {
			checksum = string(data[end+1:])
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
				DebugLog("Obj found; header:" + header + " checksum:" + checksum)
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
	// get any flags that are configured in anyproxy package
	gProxyConnectionImpl = flag.Lookup("proxyConnection").Value.(flag.Getter).Get().(string)
	gLoggingLevel = flag.Lookup("v").Value.(flag.Getter).Get().(int)
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
		destination = "/battery"
		res, _ := json.Marshal(v)
		publishMQTT(gClient, gMQTTTopic+destination, string(res))
		InfoLog(fmt.Sprint("MQTT Published::", v))
	case ConfigRS:
		gLastConfigRS = v
		destination = gAttributesTopic
		res, _ := json.Marshal(v)
		publishMQTT(gClient, gMQTTTopic+destination, string(res))
		InfoLog(fmt.Sprint("MQTT Published::", v))
	default:
		DebugLog(fmt.Sprint("Ignored message type::", v))
	}
}
