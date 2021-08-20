package alphaess

import (
	anyproxy "github.com/230delphi/go-any-proxy/anyproxy"
	log "github.com/zdannar/flogger"
	"io"
	"net"
	"os"
)

//	"rk-go-any-proxy/anyproxy/anyproxy"

type AlphaESSProxyConnection struct {
	name string
}

func (into *AlphaESSProxyConnection) copyProxyConnection(dst io.ReadWriteCloser, src io.ReadWriteCloser, dstname string, srcname string) {
	if dst == nil {
		log.Debugf("copy(): oops, dst is nil!")
		return
	}
	if src == nil {
		log.Debugf("copy(): oops, src is nil!")
		return
	}
	var err error
	// RK duplicate stream
	myfilename := getUniqueFilename(srcname)
	log.Debugf("writing file", myfilename)
	f, err := os.Create(myfilename)
	check(err)
	var buf2 io.ReadWriteCloser
	buf2 = io.ReadWriteCloser(f)
	output := io.MultiWriter(dst, buf2)
	_, err = io.Copy(output, src)
	err2 := buf2.Close()
	check(err2)
	if err != nil {
		if operr, ok := err.(*net.OpError); ok {
			if srcname == "directserver" || srcname == "proxyserver" {
				log.Debugf("copy(): %s->%s: Op=%s, Net=%s, Addr=%v, Err=%v", srcname, dstname, operr.Op, operr.Net, operr.Addr, operr.Err)
			}
			if operr.Op == "read" {
				if srcname == "proxyserver" {
					anyproxy.incrProxyServerReadErr()
				}
				if srcname == "directserver" {
					anyproxy.incrDirectServerReadErr()
				}
			}
			if operr.Op == "write" {
				if srcname == "proxyserver" {
					anyproxy.incrProxyServerWriteErr()
				}
				if srcname == "directserver" {
					anyproxy.incrDirectServerWriteErr()
				}
			}
		}
	}
	dst.Close()
	src.Close()
}
