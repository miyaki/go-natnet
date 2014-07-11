package main

import (
	. "bitbucket.org/miyaki/go-natnet"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	//osc stream
)

func main() {
	pdata := make([]byte, BUFFER_SIZE)
	f, err := os.Open("payload.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	n, err := f.Read(pdata)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data: ", n)
	//r := bufio.NewReader(bytes.NewBuffer(pdata))
	r := bytes.NewBuffer(pdata)

	p := NewPacket()
	p.Decode(r)
	fmt.Println(p)

	fmt.Printf("natnet client\n")

	////////////////////////////////////////////////////////////////////////////////////
	//client := newNatNetClient("127.0.0.1", 9999)
	//client.addHandler()

	//client.ListenAndDispatch()

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprint(MCAST_GRP, ":", DATA_PORT))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// todo make it to channel
	for {
		buf := make([]byte, BUFFER_SIZE)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("server read:", n)

		p := NewPacket()
		//p.decode(bufio.NewReader(bytes.NewBuffer(buf)))
		p.Decode(bytes.NewBuffer(buf))
		fmt.Println(p)
	}
}
