package alphaess

import (
	"encoding/json"
	"fmt"
	"github.com/230delphi/go-any-proxy/anyproxy"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namsral/flag"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

const MinMessageSize = 10
const COMMANDRQTOPIC = "/LastCommand"
const SERIALRQTOPIC = "/LastSerialRQ"
const ATTRIBUTESTOPIC = "/attributes"
const CHARGEBATTERYTOPIC = "/action/chargebattery"

var COMMANDHEADER = []byte{0x1, 0x1, 0x8}
var SUCCESSHEADER1 = []byte{0x1, 0x2, 0x0}
var SUCCESSHEADER2 = []byte{0x1, 0x2, 0x4}
var CONFIGHEADER = []byte{0x1, 0x2, 0x3}

var gClient mqtt.Client
var gMQTTBrokerAddress string
var gMQTTUser string
var gMQTTPassword string
var gAlphaEssInstance string
var gLogList string
var gTZLocation string
var gMQTTTimeoutSeconds int
var gTopicBase string
var gMQTTTopic string
var gProxyConnectionImpl string
var gMQTTBase = "homeassistant"
var gChargeBatteryTopic string
var gCommandRQTopic string

var gLastServerConfig ConfigRS
var gLastClientConfig ConfigRS
var gLastChargeState = false
var gLastCommandRQ CommandRQ
var gLastSerialRQ SerialRQ

var gActiveConversations []*conversationType = nil
var gSystemStarted = false

// TODO integration tests for alphaESS MQTT config and implementation.

type ChargeBatteryAction struct {
	GridCharge      bool `json:"GridCharge,string"`         // set to charge or not
	StartHour       int  `json:"StartHour,omitempty"`       // optional hour to start
	MinimumDuration int  `json:"MinimumDuration,omitempty"` // minimum minutes to be charging
	BatHighCap      int  `json:"BatHighCap,omitempty"`      // when to stop charging
}

func initFlagConfig() {
	flag.StringVar(&gMQTTBrokerAddress, "MQTTAddress", "tcp://127.0.0.1:1883", "MQTT address. Example: tcp://127.0.0.1:1883\n")
	flag.StringVar(&gMQTTUser, "MQTTUser", "", "MQTT username\n")
	flag.StringVar(&gMQTTPassword, "MQTTPassword", "", "MQTT password\n")
	flag.IntVar(&gMQTTTimeoutSeconds, "MQTTSendTimeout", 5, "MQTT timeout for sending message\n")
	flag.StringVar(&gTopicBase, "MQTTTopicBase", gMQTTBase+"/sensor/", "MQTT base topic. ")
	flag.StringVar(&gAlphaEssInstance, "AlphaESSID", "alphaess1", "AlphaESS instance name, appended to MQTTTopicBase. All data is set on this topic eg: homeassisant/sensor/alphaess1/config\n")
	flag.StringVar(&gLogList, "MSGLogging", "", "Messages to Log. Leave unset for no logging. Log all:\"*\"; log selected: \"GenericRQ,CommandIndexRQ,CommandRQ,ConfigRS,StatusRQ\"")
	flag.StringVar(&gTZLocation, "TZLocation", "Local", "Timezone override to ensure time of collection is accurate.")
	gMQTTTopic = gTopicBase + gAlphaEssInstance
	DebugLog("initFlagConfig complete" + gMQTTBrokerAddress)
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
	gChargeBatteryTopic = gMQTTTopic + CHARGEBATTERYTOPIC
	gCommandRQTopic = gMQTTTopic + COMMANDRQTOPIC
	return myClient
}

func getTimeNowInTimeZone() time.Time {
	var t time.Time
	loc, err := time.LoadLocation(gTZLocation)
	if err == nil {
		t = time.Now().In(loc)
	} else {
		t = time.Now()
	}
	return t
}

var mqtPublishHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	DebugLog("TOPIC: " + msg.Topic() + " : " + string(msg.Payload()))
	if !gSystemStarted {
		WarningLog("Ignoring MQTT Message Received on topic before system started. " + msg.Topic())
	} else {
		switch msg.Topic() {
		case gChargeBatteryTopic:
			InfoLog("MQTT Received on topic " + msg.Topic())
			action := getChargeBatteryAction(msg.Payload())
			injectConfigObject := buildConfigObject(action)
			gLastChargeState = action.GridCharge
			injectConfigBytes, _ := json.Marshal(injectConfigObject)
			injectConfigBytes = AddHeaderAndCheckSum(injectConfigBytes, CONFIGHEADER)
			if gLastCommandRQ.CmdIndex == 0 {
				gLastCommandRQ.CmdIndex = 51422064
				gLastCommandRQ.Command = "SetConfig"
			}
			gLastCommandRQ.CmdIndex = gLastCommandRQ.CmdIndex + 1
			injectSetCommand, _ := json.Marshal(gLastCommandRQ)
			injectSetCommand = AddHeaderAndCheckSum(injectSetCommand, COMMANDHEADER)
			injectSuccess1 := AddHeaderAndCheckSum([]byte("{\"Status\":\"Success\"}"), SUCCESSHEADER1)
			injectSuccess2 := AddHeaderAndCheckSum([]byte("{\"Status\":\"Success\"}"), SUCCESSHEADER2)

			configConversation := conversationType{
				indexOfNextAction: 0,
				actions: []actionType{
					{
						"start convo",
						INJECT,
						SERVER,
						injectSetCommand,
						nil,
					},
					{
						"ACK response",
						RESPOND,
						CLIENT,
						[]byte("\"Status\":\"Success\"}"),
						injectSuccess1,
					},
					{
						"RS SN w/Config",
						RESPOND,
						CLIENT,
						[]byte(SERIALRQPATTERN),
						injectConfigBytes,
					},
					{
						"ACK config RS",
						RESPOND,
						CLIENT,
						[]byte(CONFIGRSPATTERN),
						injectSuccess2,
					},
				},
				lastUpdate: 0,
			}
			gActiveConversations = []*conversationType{&configConversation}
		case COMMANDRQTOPIC:
			DebugLog("MQTT Received on topic " + msg.Topic())
			if gLastCommandRQ.CmdIndex < 10 {
				_ = json.Unmarshal(msg.Payload(), &gLastCommandRQ)
			}
		default:
			InfoLog("Ignored message Received on topic::" + msg.Topic())
		}
	}
}

func buildConfigObject(action *ChargeBatteryAction) *ConfigRS {
	var endHour = 0
	t := getTimeNowInTimeZone()
	currentHour := t.Hour()
	if action.MinimumDuration < 10 {
		action.MinimumDuration = 30
	}
	// nonzero start hour for immediate start. 0 is midnight.
	if (action.StartHour < 0) || (action.StartHour == currentHour) {
		action.StartHour = currentHour
		minMinutesThisHour := action.MinimumDuration % 60
		minHours := ((action.MinimumDuration - minMinutesThisHour) / 60) + 1
		if t.Minute() < minMinutesThisHour {
			endHour = (currentHour + minHours) % 24
		} else {
			endHour = (currentHour + minHours + 1) % 24
		}
	} else {
		minHours := int(math.Ceil(float64(1 + (action.MinimumDuration / 60))))
		endHour = (action.StartHour + minHours) % 24
	}

	if !action.GridCharge {
		endHour = action.StartHour
	}

	var injectConfigObject = ConfigRS{}
	// populate from saved values
	injectConfigObject = gLastServerConfig

	injectConfigObject.TimeChaF2 = action.StartHour
	injectConfigObject.TimeChaE2 = endHour
	injectConfigObject.Status = "Success"
	injectConfigObject.BackUpBox = false
	injectConfigObject.BatHighCap = float32(action.BatHighCap)
	injectConfigObject.BatHighCapWE = 0
	injectConfigObject.BatReady = 0
	injectConfigObject.BatUseCap = 10
	injectConfigObject.BatUseCapWE = 0
	injectConfigObject.Generator = false
	injectConfigObject.GridCharge = action.GridCharge
	injectConfigObject.GridChargeWE = false
	injectConfigObject.CtrDis = false
	injectConfigObject.CtrDisWE = false
	injectConfigObject.SelfUseOrEconomic = 0
	injectConfigObject.ReliefMode = 0

	return &injectConfigObject
}

func getChargeBatteryAction(payload []byte) (batteryAction *ChargeBatteryAction) {
	action := ChargeBatteryAction{!gLastChargeState, -1, 10, 50}
	err := json.Unmarshal(payload, &action)
	if err != nil {
		ExceptionLog(err, "getChargeBatteryAction() Error Unmarshalling MQTT Payload: "+string(payload))
	}
	DebugLog(fmt.Sprintf("ChargeBatteryAction: %v \n from: %s", action, string(payload)))
	return &action
}

func connLostHandler(client mqtt.Client, err error) {
	ErrorLog(fmt.Sprintf("Connection lost, reason: %v", err))
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
			ErrorLog(fmt.Sprint("Failed to send message: %v\n", t.Error()))
			ExceptionLog(t.Error(), "publishMQTT()")
		}
		if !success {
			ErrorLog("Timeout occurred after: " + strconv.Itoa(gMQTTTimeoutSeconds) + " sending:" + msg)
			ExceptionLog(t.Error(), "publishMQTT()")
		}
	}()
}

func subscribeTopic(mqClient mqtt.Client, topic string) {
	token := mqClient.Subscribe(topic, 1, nil)
	token.Wait()
	DebugLog("Subscribed to topic: " + topic)
}

func init() {
	initFlagConfig()
	anyproxy.InitConfig()
	// get any flags that are configured in anyproxy package
	gProxyConnectionImpl = flag.Lookup("proxyConnection").Value.(flag.Getter).Get().(string)
	gLoggingLevel = flag.Lookup("v").Value.(flag.Getter).Get().(int)
	gClient = initMQTT()
}

func publishAlphaESSStats(obj Response, source string) (result bool) {
	var destination string
	logResponseObject(obj, source)
	result = true
	switch v := obj.(type) {
	case StatusRQ:
		destination = "/state"
		res, _ := json.Marshal(v)
		publishMQTT(gClient, gMQTTTopic+destination, string(res))
		InfoLog(fmt.Sprint("SRC:"+source+"; MQTT Published::", v))
	case BatteryRQ:
		destination = "/battery"
		res, _ := json.Marshal(v)
		publishMQTT(gClient, gMQTTTopic+destination, string(res))
		InfoLog(fmt.Sprint("SRC:"+source+"; MQTT Published::", v))
	case ConfigRS:
		destination = ATTRIBUTESTOPIC
		res, _ := json.Marshal(v)
		publishMQTT(gClient, gMQTTTopic+destination, string(res))
		if source == "directserver" {
			gLastServerConfig = v
		} else {
			gLastClientConfig = v
		}
		InfoLog(fmt.Sprint("SRC:"+source+"; MQTT Published::", v))
	case CommandRQ:
		gLastCommandRQ = v
		destination = COMMANDRQTOPIC
		res, _ := json.Marshal(v)
		publishMQTT(gClient, gMQTTTopic+destination, string(res))
		InfoLog(fmt.Sprint("SRC:"+source+"; MQTT Published::", v))
	case SerialRQ:
		gLastSerialRQ = v
		destination = SERIALRQTOPIC
		res, _ := json.Marshal(v)
		publishMQTT(gClient, gMQTTTopic+destination, string(res))
		InfoLog(fmt.Sprint("SRC:"+source+"; MQTT Published::", v))
	default:
		DebugLog(fmt.Sprint("SRC:"+source+"; MQTT not publishing message type::", v))
		result = false
	}
	return result
}
