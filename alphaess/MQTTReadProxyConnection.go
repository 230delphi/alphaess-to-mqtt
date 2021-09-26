package alphaess

import (
	"github.com/230delphi/go-any-proxy/anyproxy"
	log "github.com/zdannar/flogger"
	"io"
	"net"
	"os"
)

type MQTTReadProxyConnection struct {
	LogToFile bool
}

func (into *MQTTReadProxyConnection) setLogToFile() {
	into.LogToFile = true
}

func (into *MQTTReadProxyConnection) SpawnBiDirectionalCopy(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
	//ReadOnly, simply extract for MQTT and copy
	go into.CopyProxyConnection(dst, src, dstName, srcName)
	go into.CopyProxyConnection(src, dst, srcName, dstName)
}

func (into *MQTTReadProxyConnection) CopyProxyConnection(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
	if dst == nil {
		log.Debugf("copy(): oops, dst is nil!")
		return
	}
	if src == nil {
		log.Debugf("copy(): oops, src is nil!")
		return
	}
	var err error
	var output io.Writer
	var mqttDst io.Writer = &MQTTWriter{mySource: srcName}
	var buf2 io.ReadWriteCloser

	if into.LogToFile {
		//TODO retest this logging mechanism
		myFilename := getUniqueFilename(srcName)
		log.Debugf("writing file", myFilename)
		f, err := os.Create(myFilename)
		CheckError(err)
		buf2 = io.ReadWriteCloser(f)
		output = io.MultiWriter(dst, mqttDst, buf2)
	} else {
		output = io.MultiWriter(dst, mqttDst)
	}
	_, err = io.Copy(output, src)

	ExceptionLog(err)
	if buf2 != nil {
		_ = buf2.Close()
	}

	ReportStatistics(err, srcName, dstName)
	err = dst.Close()
	ExceptionLog(err)
	err = src.Close()
	ExceptionLog(err)
}

func ReportStatistics(err error, srcName string, dstName string) {
	//TODO move to function in anyproxy
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {
			if srcName == "directserver" || srcName == "proxyserver" {
				log.Debugf("copy(): %s->%s: Op=%s, Net=%s, Addr=%v, Err=%v", srcName, dstName, opError.Op, opError.Net, opError.Addr, opError.Err)
			}
			if opError.Op == "read" {
				if srcName == "proxyserver" {
					anyproxy.IncrProxyServerReadErr()
				}
				if srcName == "directserver" {
					anyproxy.IncrDirectServerReadErr()
				}
			}
			if opError.Op == "write" {
				if srcName == "proxyserver" {
					anyproxy.IncrProxyServerWriteErr()
				}
				if srcName == "directserver" {
					anyproxy.IncrDirectServerWriteErr()
				}
			}
		}
	}
}
