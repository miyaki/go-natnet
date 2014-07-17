package natnet

import (
	"errors"
	"net"
)

// NatNet network settings
var (
	SERVER_IP    = "127.0.0.1"
	MCAST_GRP    = "239.255.42.99"
	COMMAND_PORT = 1510
	DATA_PORT    = 1511
	BUFFER_SIZE  = 1024
	VERSION      = 0.0
)

// NatNetClient protocol receiver
type NatNetClient struct {
	Address    string
	Port       int
	Dispatcher *Dispatcher
	running    bool
	conn       *net.UDPConn
}

// TODO for alt interface to bind multicast socket

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
func (client *NatNetClient) ListenAndDispatch() error {
	if client.running {
		return errors.New("Client is already listning")
	}
	return nil
}
