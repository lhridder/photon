package photon

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type Gateway struct {
	Listener  net.Listener
	ListenTo  string
	Proxies   map[string]Proxy
	conngroup sync.WaitGroup
}

func (gateway *Gateway) Listen() error {
	listener, err := net.Listen("udp", gateway.ListenTo)
	if err != nil {
		return fmt.Errorf("failed to create listener: %s", err)
	}

	gateway.Listener = listener
	return nil
}

func (gateway *Gateway) Serve() {
	for {
		conn, err := gateway.Listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("Closing listener on", gateway.ListenTo)
				//TODO register closed listener
				return
			}
			continue
		}

		go func() {
			gateway.conngroup.Add(1)
		}()
	}

}
