package alphaess

import (
	"github.com/230delphi/go-any-proxy/anyproxy"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type MQTTStubProxyConnection struct {
	LogToFile bool
}

func (into *MQTTStubProxyConnection) setLogToFile() {
	into.LogToFile = true
}

func (into *MQTTStubProxyConnection) SpawnBiDirectionalCopy(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
	//TODO either server impl, or switch from proxy
	//go into.CopyProxyConnection(dst, src, dstName, srcName)
	//go into.CopyProxyConnection(src, dst, srcName, dstName)
}

func (into *MQTTStubProxyConnection) CopyProxyConnection(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
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
	// RK duplicate stream

	if into.LogToFile {
		//TODO test this logging mechanism - needed?
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

	ExceptionLog(err, "SRC:"+srcName)
	if buf2 != nil {
		_ = buf2.Close()
	}

	anyproxy.ReportStatistics(err, srcName, dstName)
	err = dst.Close()
	ExceptionLog(err, "SRC:"+srcName)
	err = src.Close()
	ExceptionLog(err, "SRC:"+srcName)
}
