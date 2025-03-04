package internal

import (
	"github.com/netbirdio/netbird/client/internal/dns"
	"github.com/netbirdio/netbird/client/internal/listener"
	"github.com/netbirdio/netbird/client/internal/stdnet"
	"github.com/netbirdio/netbird/iface"
)

// MobileDependency collect all dependencies for mobile platform
type MobileDependency struct {
	TunAdapter            iface.TunAdapter
	IFaceDiscover         stdnet.ExternalIFaceDiscover
	NetworkChangeListener listener.NetworkChangeListener
	HostDNSAddresses      []string
	DnsReadyListener      dns.ReadyListener
	DnsManager            dns.IosDnsManager
	FileDescriptor        int32
}
