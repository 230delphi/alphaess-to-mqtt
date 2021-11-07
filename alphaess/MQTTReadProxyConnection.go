package alphaess

import (
	"github.com/230delphi/go-any-proxy/anyproxy"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

var LogToFile bool

type MQTTReadProxyConnection struct {
}

func (into *MQTTReadProxyConnection) setLogToFile() {
	LogToFile = true
}

func (into *MQTTReadProxyConnection) SpawnBiDirectionalCopy(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
	//ReadOnly, simply extract for MQTT and copy
	go into.CopyProxyConnection(dst, src, dstName, srcName)
	go into.CopyProxyConnection(src, dst, srcName, dstName)
}
func (into *MQTTReadProxyConnection) CopyProxyConnection(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
	SimpleCopyProxyConnection(dst, src, dstName, srcName)
}

func SimpleCopyProxyConnection(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
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

	if LogToFile {
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
