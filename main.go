package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	//hid stream
	//osc stream
	//usb serial stream
)

// NatNet network settings
var (
	SERVER_IP    = "127.0.0.1"
	MCAST_GRP    = "239.255.42.99"
	COMMAND_PORT = 1510
	DATA_PORT    = 1511
	BUFFER_SIZE  = 1024
)

// NatNetClient protocol receiver
type NatNetClient struct {
	Address    string
	Port       int
	Dispatcher *Dispatcher
	running    bool
	conn       *net.UDPConn
}

// Dispatcher dispatches message
type Dispatcher struct {
	handlers map[string]Handler
}

func newDispatcher() (dispatcher *Dispatcher) {
	return &Dispatcher{handlers: make(map[string]Handler)}
}

// Handler for messages
type Handler interface {
	//	HandleMessage(msg *NatNetMessage)
}

//type HandlerFunc func(msg *NatNetMessage)

func newNatNetClient(address string, port int) (client *NatNetClient) {
	return &NatNetClient{
		Address:    address,
		Port:       port,
		Dispatcher: newDispatcher(),
	}
}

//Close closes NatNet connection
func (client *NatNetClient) Close() error {
	client.running = false
	return client.conn.Close()
}

//func (self *Dispatcher) Dispatch(packet NatNetPacket) {
//	fmt.Printf("natnet packet")
//}

//func (self *NatNetClient) addHandler(address string, handler HandlerFunc) error {
//return self.Dispatcher.AddMsgHandler(address, handler)
//}

// ListenAndDispatch handles packet
func (self *NatNetClient) ListenAndDispatch() error {
	if self.running {
		return errors.New("Client is already listning")
	}
	return nil
}

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
	p.decode(r)
	fmt.Println(p)

	fmt.Printf("natnet client\n")
	//client := newNatNetClient("127.0.0.1", 9999)
	//client.addHandler()

	//client.ListenAndDispatch()

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprint(MCAST_GRP, ":", DATA_PORT))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		buf := make([]byte, BUFFER_SIZE)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("server read:", n)

		p := NewPacket()
		//p.decode(bufio.NewReader(bytes.NewBuffer(buf)))
		p.decode(bytes.NewBuffer(buf))
		fmt.Println(p)
	}
}
