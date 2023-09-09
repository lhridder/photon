package photon

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lhridder/photon/config"
	"github.com/lhridder/photon/protocol"
	"log"
	"net"
	"sync"
)

type Gateway struct {
	Listener  net.Listener
	ListenTo  string
	Proxies   map[string]Proxy
	conngroup sync.WaitGroup
	Cfg       config.Globalconfig
}

func (gw *Gateway) Listen() error {
	listener, err := net.Listen("tcp", gw.ListenTo)
	if err != nil {
		return fmt.Errorf("failed to create listener: %s", err)
	}

	gw.Listener = listener
	return nil
}

func (gw *Gateway) Serve() {
	for {
		conn, err := gw.Listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("Closing listener on", gw.ListenTo)
				//TODO register closed listener
				return
			}
			continue
		}

		go func() {
			log.Printf("New connection: %s", conn.RemoteAddr())
			gw.conngroup.Add(1)
			defer gw.conngroup.Done()

			mc := protocol.Conn{NetConn: conn}

			err = gw.handleConnection(&mc)
			if err != nil {
				log.Println(err)
			}

			mc.NetConn.Close()
		}()
	}

}

func (gw *Gateway) handleConnection(mc *protocol.Conn) error {
	pk, err := mc.ReadPacket(1024)
	if err != nil {
		return err
	}

	hs, err := pk.ReadHandshake()
	if err != nil {
		return err
	}

	switch hs.NextState {
	case 1:
		mc.Type = protocol.Status
	case 2:
		mc.Type = protocol.Login
	default:
		//TODO mark invalid
	}

	proxy, ok := gw.Proxies[hs.ServerAddr]
	if !ok {
		err = gw.handleUnknown(mc)
		if err != nil {
			return err
		}
	}

	switch mc.Type {
	case protocol.Status:
		err = proxy.handleStatusRequest(mc)
		if err != nil {
			return err
		}
	case protocol.Login:
		err = proxy.handleLoginRequest(mc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gw *Gateway) handleUnknown(mc *protocol.Conn) error {
	if mc.Type == protocol.Status {
		pk, err := mc.ReadPacket(2)
		if err != nil {
			return err
		}

		err = pk.ReadStatusRequest()
		if err != nil {
			return err
		}

		err = gw.sendDefaultStatus(mc)
		if err != nil {
			return err
		}
	}
	if mc.Type == protocol.Login {

	}

	return nil
}

func (gw *Gateway) sendDefaultStatus(mc *protocol.Conn) error {
	response, err := json.Marshal(gw.Cfg.DefaultStatusResponse)
	if err != nil {
		return protocol.ErrJsonMarshal
	}

	spk, err := protocol.WriteStatusResponse(response)
	if err != nil {
		return err
	}

	log.Println(string(*spk.Data))

	err = mc.WritePacket(spk)
	if err != nil {
		return err
	}

	pk, err := mc.ReadPacket(9)
	if err != nil {
		return err
	}

	ping, err := pk.ReadPingRequest()
	if err != nil {
		return err
	}

	ppk, err := protocol.WritePingResponse(ping)
	if err != nil {
		return err
	}

	err = mc.WritePacket(ppk)
	if err != nil {
		return err
	}

	return nil
}

func (gw *Gateway) sendDefaultLogin(mc *protocol.Conn) error {
	//TODO implement

	return nil
}
