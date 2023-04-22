package photon

type Proxy struct {
	UID           string
	ProxyTo       string
	ProxyProtocol bool
	Domainnames   []string
}
