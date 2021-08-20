package main

import (
	anyproxy "github.com/230delphi/go-any-proxy/anyproxy"
	"bufio"
	"main/alphaess"
	"os"
)

var myProxyConnection anyproxy.ProxyConnection

func test() {
	myfilename := "client_src.stream"
	debugLog("test file:" + myfilename)
	f, err := os.Open(myfilename)
	var myreader = bufio.NewReader(f)
	check(err)
	alphaess.publishHASEntityConfig()
	processStream(err, myreader)
	//	TODO complete tests
}

func main() {
	myProxyConnection = anyproxy.ProxyConnection & AlphaESSProxyConnection{}
	alphaess.PublishHASEntityConfig()
	//	processStream(err, myreader)
	anyproxy.StartProxy(myProxyConnection)
}
