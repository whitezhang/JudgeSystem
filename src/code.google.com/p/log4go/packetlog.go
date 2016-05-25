/* packetlog.go - packet-oriented log writer */
/*
modification history
--------------------
2015/01/05, by taochunhua, create
*/
/*
DESCRIPTION
Packet-oriented log writer. Using PacketConn interface in golang (net.PacketConn).
The network net must be a packet-oriented network: 
udp, udp4, udp6, unixgram
*/

package log4go

import (
    "errors"
	"fmt"
	"net"
	"os"
    "strings"
)

// packet connection
type PacketConn struct {
    conn            net.PacketConn
    remoteAddr      net.Addr
}

var (
    ErrNetworkNotMatch = errors.New("network in following list:  udp, udp4, udp6, unixgram")
)

func resolveAddr(network string, remoteAddr string) (net.Addr, error) {
    var addr net.Addr
    var err error
    
    network = strings.ToLower(network)
    switch network {
    case "udp":  fallthrough
    case "udp4": fallthrough
    case "udp6": addr, err = net.ResolveUDPAddr(network, remoteAddr)
    case "unixgram": addr, err = net.ResolveUnixAddr(network, remoteAddr)
    default: addr = nil
             err = ErrNetworkNotMatch
    }
    
    return addr, err
}

// create packet connection
func newPacketConn(network string, remoteAddr string) (*PacketConn, error) {
    network = strings.ToLower(network)
    if network != "udp" &&  network != "udp4" && 
       network != "udp6" &&  network != "unixgram" {
        return nil, ErrNetworkNotMatch
    }
    // create golang net.PacketConn
    // Do not receive data, so local address is set to "".
    conn, err := net.ListenPacket(network, "")
    if err != nil {
        return nil, err
    }
    
    // resolve address
    address, err := resolveAddr(network, remoteAddr)
    if err != nil {
        return nil, err
    }
    
    // create PacketConn
    pc := &PacketConn{conn:conn, remoteAddr:address}
    
    return pc, nil
}

// send data to log server
func (pc *PacketConn) Send(data []byte) error {
    _, err := pc.conn.WriteTo(data, pc.remoteAddr)
    return err
}

// This log writer sends output to a packet connection
type PacketWriter struct {
    LogCloser   //for Elegant exit
    
    rec     chan *LogRecord
    conn    *PacketConn
    name    string
}

// send data
func (w *PacketWriter) Send(data []byte) error {
    return w.conn.Send(data)
}

// This is the PacketWriter's output method
func (w *PacketWriter) LogWrite(rec *LogRecord) {
    if !LogWithBlocking {
        if len(w.rec) >= LogBufferLength {
            if WithModuleState {
                log4goState.Inc("ERR_SOCK_LOG_OVERFLOW", 1)
            }            
            
            return
        }
    }
    
	w.rec <- rec
}

// get writer name
func (w *PacketWriter) Name() string {
    return w.name
}

// get rec channel length
func (w *PacketWriter) QueueLen() int {
    return len(w.rec)
}

func NewPacketWriter(name string, network string, 
                     remoteAddr string, format string) *PacketWriter {
	conn, err := newPacketConn(network, remoteAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewPacketWriter(%s, %s): %s\n", 
                    name, remoteAddr, err)
		return nil
	}
    
	w := &PacketWriter{
		rec:      make(chan *LogRecord, LogBufferLength),
        conn:     conn,
		name:     name,
	}
    
    //init LogCloser
    w.LogCloserInit()
    
    // add w to collection of all writers' info
    writersInfo = append(writersInfo, w)

	go func() {
        for {
            rec := <-w.rec
            
            if w.EndNotify(rec) {
                return
            }
            
            if rec.Binary != nil {
                w.Send(rec.Binary)
                putBuffer(rec.Binary) // Binary is allocated from buffer pool
            } else {
                msg := FormatLogRecord(format, rec)
                w.Send([]byte(msg))
            }
		}
	}()

	return w
}

//wait for dump all log and close chan
func (w *PacketWriter) Close() {
	w.WaitForEnd(w.rec)
    close(w.rec)
}
