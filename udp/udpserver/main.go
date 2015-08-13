package main

import (
	"flag"
	"fmt"
	"net"
	"log"
)

var port = flag.Int("p", 1053, "The listening port of the udp server")
var udpPackageBufferSize = flag.Int("l", 100, "The size of the udp package buffer")

func main() {
	flag.Parse()

	// 创建监听
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: *port,
	})

	if err != nil {
		fmt.Println("listen failed!", err)
		return
	}
	defer socket.Close()

	for {
		// 读取数据
		data := make([]byte, *udpPackageBufferSize)
		readn, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("recvfrom error!", err)
			continue
		}
		
		log.Println("recvfrom dlen", len(data))

		go process(socket, data[:readn], remoteAddr)
	}
}

func process(conn *net.UDPConn, data []byte, remoteAddr *net.UDPAddr) {
	_, err := conn.WriteToUDP(data, remoteAddr)
	if err != nil {
		fmt.Println("send data error", err)
	}
}
