package alphaess

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		debugLog(e.Error())
		panic(e)
	}
}
func debugLog(msg string) {
	if false {
		log.Println(msg)
	}
}

func errorLog(msg string) {
	log.Println("ERROR:" + msg)
}
func expLog(emsg error) {
	errorLog("EXP:" + emsg.Error())
}

func initMQTT() (myClient mqtt.Client) {
	debugLog("init MQTT")
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	//mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	opts := mqtt.NewClientOptions().AddBroker(mqttBrokerAddress)
	opts.SetClientID("alphaess-gclient")
	opts.SetUsername(mqttUser)
	opts.SetPassword(mqttPassword)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	myClient = mqtt.NewClient(opts)
	if token := myClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if myClient.IsConnectionOpen() {
		debugLog("initMQTT: connected")
	}
	return myClient
}

func publishMQTT(topic string, msg string) {
	var qos byte
	var retained = false

	t := gclient.Publish(topic, qos, retained, msg)
	go func() {
		_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
		if t.Error() != nil {
			expLog(t.Error()) // Use your preferred logging technique (or just fmt.Printf)
		}
	}()
}
