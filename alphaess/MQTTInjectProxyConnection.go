package alphaess

import (
	"errors"
	"github.com/230delphi/go-any-proxy/anyproxy"
	"io"
	"os"
	"sync"
)

// TODO This & Filters could be abstracted as a more generic mechanism in the core any_proxy module.

// MQTTInjectProxyConnection is an implementation of ProxyConnection allowing Filters to process conversations
// across a connection. In addition, data is extracted and exposed via MQTT.
type MQTTInjectProxyConnection struct {
}

func (into *MQTTInjectProxyConnection) setLogToFile() {
	LogToFile = true
}

func (into *MQTTInjectProxyConnection) SpawnBiDirectionalCopy(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
	WarningLog("MQTTInjectProxyConnection: CopyAndInjectProxyConnection() is in testing only!")
	initMQTTClient()
	myServerFilter := ServerFilter{"Server"}
	myClientFilter := ClientFilter{"Client"}
	go into.CopyAndInjectProxyConnection(dst, src, dstName, srcName, &myServerFilter)
	go into.CopyAndInjectProxyConnection(src, dst, srcName, dstName, &myClientFilter)
}

func (into *MQTTInjectProxyConnection) CopyProxyConnection(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string) {
	SimpleCopyProxyConnection(dst, src, dstName, srcName)
}

func initMQTTClient() {
	if gClient == nil {
		gClient = initMQTT()
	}
	subscribeTopic(gClient, gChargeBatteryTopic)
	subscribeTopic(gClient, gCommandRQTopic)
}

//mutex is required to ensure writes can happen both directions between the 2 threads
var mutex = new(sync.RWMutex)

func (into *MQTTInjectProxyConnection) CopyAndInjectProxyConnection(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstName string, srcName string, myFilter MessageFilter) {
	if myFilter == nil {
		WarningLog("No filter passed, using PassFilter")
		myFilter = &PassFilter{}
	}
	if dst == nil {
		DebugLog("copy(): oops, dst is nil!")
		return
	}
	if src == nil {
		DebugLog("copy(): oops, src is nil!")
		return
	}
	var err error
	var output io.Writer
	var mqttDst io.Writer = &MQTTWriter{mySource: srcName}
	var buf2 io.ReadWriteCloser

	if LogToFile {
		//TODO test this logging mechanism - needed?
		myFilename := getUniqueFilename(srcName)
		DebugLog("writing file: ", myFilename)
		f, err := os.Create(myFilename)
		CheckError(err)
		buf2 = io.ReadWriteCloser(f)
		output = io.MultiWriter(dst, mqttDst, buf2)
	} else {
		// TODO probably remove this double write and just include below.
		output = io.MultiWriter(dst, mqttDst)
	}
	_, err = into.copyBufferAndFilter(output, src, nil, myFilter)

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

// for some reason this is not exported from package IO
// errInvalidWrite means that a 'write' returned an impossible count.
var errInvalidWrite = errors.New("invalid write result")

// taken from: https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/io/io.go;drc=refs%2Ftags%2Fgo1.17.1;l=402
// copyBufferAndFilter is a modified version of the core: copyBuffer (from IO) which is the underlying implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
// MessageFilter will be applied to data read from src stream. Based on the filter response it could:
//		a. pass the original data to the dst
//  	b. modify and pass to dst.
//   	c. drop the data
//		d. only respond back to the src with a message
//		and/or inject new data to the stream after specific bytes were sent
func (into *MQTTInjectProxyConnection) copyBufferAndFilter(dst io.Writer, srcReadWriter io.ReadWriteCloser, buf []byte, filter MessageFilter) (written int64, err error) {
	// map to Reader for most method.
	var src io.Reader = srcReadWriter
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
		size := 32 * 1024 * 3 // triple size of the buffer, so we always get them in one read.
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
			msgBody, valid := parseAndDebugMessage("FROM:"+filter.getName(), buf[0:nr])
			// Apply filter to message received.
			newBuffer, newResponse := filter.FilterMessages(buf, nr)
			if newBuffer != nil && valid { // Forward new buffer including any possible changes
				buf = newBuffer
				mutex.Lock() // mutex is now required because we could be writing to either stream from each thread.
				nw, ew := dst.Write(buf[0:nr])
				mutex.Unlock()
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
			} else if newResponse != nil && valid { // DROP message, and RESPOND with new
				// publish last message per normal
				publishAlphaEssBytes(msgBody, filter.getName())
				mutex.Lock()
				nrCount := len(newResponse)
				nw, ew := srcReadWriter.Write(newResponse)
				mutex.Unlock()
				// not needed parseAndDebugMessage("RESPOND to:"+filter.getName(), newResponse)
				if nw < 0 || nrCount < nw {
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
				DebugLog("Filter:" + filter.getName() + " DROP: '" + string(buf[0:nr]) + "' and RESPOND with:" + string(newResponse))
			} else if valid {
				// DROP
				DebugLog("Filter:" + filter.getName() + " DROP: " + string(buf[0:nr]))
			}
		}
		// Filter decides if it should inject into stream.
		filter.InjectMessage(dst, buf[0:nr])
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
