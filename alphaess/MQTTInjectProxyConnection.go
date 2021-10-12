package alphaess

import (
	"errors"
	log "github.com/zdannar/flogger"
	"io"
	"os"
)

type MQTTInjectProxyConnection struct {
}

func (into *MQTTInjectProxyConnection) setLogToFile() {
	LogToFile = true
}

func (into *MQTTInjectProxyConnection) SpawnBiDirectionalCopy(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
	//TODO either server impl, or switch from proxy
	go into.CopyProxyConnection(dst, src, dstName, srcName)
	go into.CopyAndInjectProxyConnection(src, dst, srcName, dstName)
}

func (into *MQTTInjectProxyConnection) CopyProxyConnection(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
	SimpleCopyProxyConnection(dst, src, dstName, srcName)
}

func initMQTTClient() {
	if gClient == nil {
		gClient = initMQTT()
	}
	subscribeTopic(gClient, gChargeBatteryTopic)
}

func (into *MQTTInjectProxyConnection) CopyAndInjectProxyConnection(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
	initMQTTClient()
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
	_, err = into.copyBufferAndInject(output, src, nil)

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

// for some reason this is not exported from package IO
// errInvalidWrite means that a write returned an impossible count.
var errInvalidWrite = errors.New("invalid write result")

// taken from: https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/io/io.go;drc=refs%2Ftags%2Fgo1.17.1;l=402
// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func (into *MQTTInjectProxyConnection) copyBufferAndInject(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}

		//TODO I dont have the checksum - this probably needs to be added.
		if nr < 20 && ClientInject != nil && len(ClientInject) > 0 {
			clientStr := string(ClientInject)
			DebugLog("Writing to path: " + clientStr)
			myMsgBytes := []byte(ClientInject)
			writen, err := dst.Write(myMsgBytes)
			if err != nil {
				ErrorLog("Error writing message: " + clientStr)
				ExceptionLog(err)
			}
			if writen == len(ClientInject) {
				DebugLog("Inject Write successful:" + clientStr)
			} else {
				ErrorLog("Inject Write Failed:" + clientStr)
			}
			ClientInject = nil
		}

		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
