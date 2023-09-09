package photon

import "github.com/lhridder/photon/protocol"

type Proxy struct {
	UID           string
	ProxyTo       string
	ProxyProtocol bool
	Domainnames   []string
}

func (p *Proxy) handleStatusRequest(mc *protocol.Conn) error {
	//TODO implement

	return nil
}

func (p *Proxy) handleLoginRequest(mc *protocol.Conn) error {
	//TODO implement

	return nil
}
