package main

import (
	"github.com/230delphi/go-any-proxy/anyproxy"
	"main/alphaess"
)

var myProxyConnection anyproxy.ProxyConnectionManager

func main() {
	myProxyConnection = alphaess.GetMQTTConnection()
	alphaess.PublishHASEntityConfig()
	anyproxy.StartProxy(myProxyConnection)

	//var injectConfigObject = alphaess.ConfigRS{}
	//injectConfigObject.TimeChaF2 = 1
	//injectConfigObject.TimeChaE2 = 2
	//injectConfigObject.Generator = true
	//injectConfig, _ := json.Marshal(injectConfigObject)
	//fmt.Println(string(injectConfig))
	//fmt.Println("NEXT")
	//var injectConfigObject2 = alphaess.ConfigRS{}
	//injectConfigObject2 = injectConfigObject
	//injectConfigObject2.GridCharge = true
	//injectConfig2, _ := json.Marshal(injectConfigObject2)
	//fmt.Println(string(injectConfig2))
}
