package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/funny/overall"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var concurrence = flag.Int("c", 1, "Number of multiple requests to make")
var numRequest = flag.Int("n", 1, "Number of requests to perform")
var verbosity = flag.Bool("v", false, "How much troubleshooting info to print")
var echoVerify = flag.Bool("e", false, "the echo mode and verify the response data")
var hostPort = flag.String("h", "127.0.0.1:1053", "The hostname and port of the udp server")
var messageLen = flag.Int("l", 100, "The length of the send message")
var timeoutMs = flag.Int("t", 100, "The timeout milliseconds")

type Stat struct {
	send       int64
	recv       int64
	sendsucc   int64
	sendfail   int64
	recvsucc   int64
	recvfail   int64
	compareerr int64
}

func (s Stat) dump() string {
	return fmt.Sprintf("send:%v recv:%v sendsucc:%v recvsucc:%v compareerr:%v sendfail:%v recvfail:%v",
		s.send, s.recv, s.sendsucc, s.recvsucc, s.compareerr, s.sendfail, s.recvfail)
}

var stat Stat

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	recoder := overall.NewTimeRecorder()
	for c := 0; c < *concurrence; c++ {
		go request(&wg, recoder)
		wg.Add(1)
	}

	wg.Wait()

	fmt.Println(stat.dump())

	var buf bytes.Buffer
	recoder.WriteCSV(&buf)
	fmt.Println(buf.String())

	if stat.compareerr != 0 || stat.recvfail != 0 || stat.sendfail != 0 {
		os.Exit(1)
	}
}

func request(wg *sync.WaitGroup, record *overall.TimeRecorder) {
	defer wg.Done()

	addr, err := net.ResolveUDPAddr("udp", *hostPort)
	if err != nil {
		fmt.Println("server address error. It MUST be a format like this hostname:port", err)
		return
	}

	// Create a udp socket and connect to server
	socket, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		fmt.Printf("connect to udpserver %v failed : %v", addr.String(), err.Error())
		return
	}
	defer socket.Close()

	msg := make([]byte, *messageLen)
	for i := 0; i < *messageLen; i++ {
		msg[i] = 'a' + byte(i)%26
	}

	for n := 0; n < *numRequest; n++ {
		t := time.Now()

		socket.SetDeadline(t.Add(time.Duration(100 * time.Millisecond)))

		// send data to server
		_, err = socket.Write(msg)
		if err != nil {
			fmt.Println("send data error ", err)
			atomic.AddInt64(&stat.sendfail, 1)
			return
		}

		// recv data from server
		data := make([]byte, *messageLen)
		read, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("recv data error ", err)
			atomic.AddInt64(&stat.recvfail, 1)
			return
		}
		data = data[:read]
		record.Record("onetrip", time.Since(t))

		if *verbosity {
			fmt.Printf("server addr [%v], response data len:%v [%s]\n", remoteAddr, read, string(data[:read]))
		}

		if *echoVerify {
			if bytes.Compare(data, msg) != 0 {
				atomic.AddInt64(&stat.compareerr, 1)
			}
		}

		atomic.AddInt64(&stat.send, 1)
		atomic.AddInt64(&stat.recv, 1)
		atomic.AddInt64(&stat.sendsucc, 1)
		atomic.AddInt64(&stat.recvsucc, 1)
	}
}
