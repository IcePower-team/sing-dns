package dns

import (
	"context"
	"net"
	"net/netip"
	"os"
	"sync"
	"time"

	"github.com/ameshkov/dnscrypt/v2"
	"github.com/miekg/dns"
	M "github.com/sagernet/sing/common/metadata"
)

var _ Transport = (*SDNSTransport)(nil)

type SDNSTransport struct {
	client      *dnscrypt.Client
	name        string
	destination string

	mtx sync.Mutex
}

func init() {
	RegisterTransport([]string{"sdns"}, func(options TransportOptions) (Transport, error) {
		return NewSDNSTransport(options), nil
	})
}

func NewSDNSTransport(options TransportOptions) *SDNSTransport {
	return &SDNSTransport{
		client: &dnscrypt.Client{
			Net:     "udp",
			Timeout: 10 * time.Second,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return options.Dialer.DialContext(ctx, network, M.ParseSocksaddr(addr))
			},
		},
		name:        options.Name,
		destination: options.Address,
	}
}

func (t *SDNSTransport) Name() string {
	return t.name
}

func (t *SDNSTransport) Start() error {
	return nil
}

func (t *SDNSTransport) Reset() {

}

func (t *SDNSTransport) Close() error {
	return nil
}

func (t *SDNSTransport) Raw() bool {
	return true
}

func (t *SDNSTransport) Exchange(ctx context.Context, message *dns.Msg) (*dns.Msg, error) {
	resolverInfo, err := t.client.Dial(t.destination)
	if err != nil {
		return nil, err
	}
	return t.client.Exchange(message, resolverInfo)
}

func (t *SDNSTransport) Lookup(ctx context.Context, domain string, strategy DomainStrategy) ([]netip.Addr, error) {
	return nil, os.ErrInvalid
}
