package alphaess

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func CheckError(e error) {
	if e != nil {
		DebugLog(e.Error())
		panic(e)
	}
}
func DebugLog(msg string) {
	//TODO debug logging config
	if false {
		log.Println(msg)
	}
}

func InfoLog(msg string) {
	//TODO info logging config
	if false {
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

func connLostHandler(client mqtt.Client, err error) {
	fmt.Printf("Connection lost, reason: %v\n", err)
	connectClient(client)
	ErrorLog("Reconnected MQTT")
}

func initMQTT() (myClient mqtt.Client) {
	DebugLog("init MQTT")
	//mqtt.DEBUG = 	log.New(os.Stderr, "[DEBUG][MQTT]", log.Ltime)
	mqtt.WARN = log.New(os.Stderr, "[WARN] [MQTT]", log.Ltime)
	mqtt.ERROR = log.New(os.Stdout, "[ERROR][MQTT]", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] [MQTT]", 0)
	opts := mqtt.NewClientOptions().AddBroker(gMQTTBrokerAddress)
	var secs int64 = time.Now().Unix()
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

	opts.SetDefaultPublishHandler(mqttPublishHandler)
	opts.SetConnectionLostHandler(connLostHandler)
	myClient = mqtt.NewClient(opts)
	if !connectClient(myClient) {
		panic("Failed to connect to MQTT: " + gMQTTBrokerAddress)
	}
	return myClient
}

var mqttPublishHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
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

func getUniqueFilename(prepend string) string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10) + "_" + prepend + "_stream.file"
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
