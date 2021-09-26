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
}
